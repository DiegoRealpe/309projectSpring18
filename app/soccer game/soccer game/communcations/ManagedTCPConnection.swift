//
//  ManagedTCPConnection.swift
//  soccer game
//
//  Created by rtoepfer on 2/14/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SwiftSocket

//wraps a SwiftSocket TCP client in order to expose application specific functionality
//and allow for potentially changing libraries in the future
class ManagedTCPConnection{
    
    var dataHandler : ([UInt8]) -> Void
    
    private var client : TCPClient
    let port : Int32
    
    //read by dispatcher queues to determine when to stop
    var stopRunning : Bool
    
    convenience init(address : String, port : Int32){
        self.init(address : address, port : port, dataHandler : defaultDataHandler(_:))
    }
    
    init(address : String, port : Int32, dataHandler : @escaping ([UInt8]) -> Void){
        self.client = TCPClient(address: address, port: port)
        self.stopRunning = false
        self.port = port
        self.dataHandler = dataHandler
        
        print("connecting to \(address), port \(port)")
        
        client.connect(timeout: 30).logError()
        startTCPCycle()
    }
    
    
    func sendTCP(message : String){
        guard !stopRunning else{
            //todo: determint behavior for this
            return
        }
        
        self.client.send(string: message).logError()
        print("sent: \"\(message)\"")
    }
    
    func sendTCP(data : [UInt8]){
        guard !stopRunning else{
            //todo: determint behavior for this
            return
        }
        
        self.client.send(data: data).logError()
        print("sent: \"\(data)\"")
    }
    
    func sendTCP(packet : SendablePacket){
        self.sendTCP(data: packet.toByteArray())
    }
    
    func stop(){
        print("stopping tcp connection")
        stopRunning = true
    }
    
    fileprivate func startTCPCycle(){
        startSocketRealLoop()
    }
    
    
    fileprivate func respondToTCPDataSent() {
        
        let data = self.client.read(50,timeout: 100)

        if let recieved = data{
            dataHandler(recieved)
        }
    }
    
    private func startSocketRealLoop(){
        DispatchQueue.global().async( execute: {
            while !self.stopRunning{
                self.respondToTCPDataSent()
            }
        })
    }
    
}

func defaultDataHandler(_ data: [UInt8]){
    print("recieved: " + String(bytes: data, encoding: .utf8)!, "with no handler")
}

