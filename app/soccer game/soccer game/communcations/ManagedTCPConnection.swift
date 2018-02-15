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
    
    func stop(){
        print("stopping tcp connection")
        stopRunning = true
    }
    
    fileprivate func startTCPCycle(){
        tcpCycle()
    }
    
    
    fileprivate func respondToTCPDataSent() {
        if let recieved = self.client.read(50){
            print("recieved: " + String(bytes: recieved, encoding: .utf8)!)
        }
    }
    
    //starts dispatch queue that calls itself after completion
    func tcpCycle(){
        DispatchQueue.main.asyncAfter(deadline: .now() + .milliseconds(250), execute: {
            self.respondToTCPDataSent()
            
            //call tcpCycle to start another dispatch Queue
            if self.stopRunning{
                self.client.close()
            }else{
                self.tcpCycle()
            }
        })
    }
    
}
