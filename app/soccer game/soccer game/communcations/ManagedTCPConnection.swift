//
//  ManagedTCPConnection.swift
//  soccer game
//
//  Created by rtoepfer on 2/14/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SwiftSocket

class ManagedTCPConnection{
    
    var datahandler : ([UInt8]) -> Void = { data in
        print("recieved: " + String(bytes: data, encoding: .utf8)!, "with no handler")
    }
    
    var client : TCPClient
    let port : Int32
    
    //read by dispatcher queues to determine when to stop
    var stopRunning : Bool
    
    init(address : String, port : Int32){
        self.client = TCPClient(address: address, port: port)
        self.stopRunning = false
        self.port = port
        
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
    func stop(){
        print("stopping tcp connection")
        stopRunning = true
    }
    
    fileprivate func startTCPCycle(){
        tcpCycle()
    }
    
    
    fileprivate func respondToTCPDataSent() {
        
        let data = self.client.read(25,timeout: 100)

        if let recieved = data{
            datahandler(recieved)
        }
    }
    
    
    
    //starts dispatch queue that calls itself after completion
    func tcpCycle(){
        DispatchQueue.global().async( execute: {
            while !self.stopRunning{
                self.respondToTCPDataSent()
            }
        })
    }
    
}

