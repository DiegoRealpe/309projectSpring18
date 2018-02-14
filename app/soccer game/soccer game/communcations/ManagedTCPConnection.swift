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

    //read by dispatcher queues to determine when to stop
    var stopRunning = true

    //is currently running
    var isRunning : Bool

    init(address : String, port : Int32, start : Bool){
        self.client = TCPClient(address: address, port: port)

        if start {
            self.isRunning = true
            self.stopRunning = false
            client.connect(timeout: 30).logError()
            startTCPCycle(client: self.client)
        }else{
            self.isRunning = false
        }
    }

    func startTCPCycle(client tcp : TCPClient){

    }


    fileprivate func respondToTCPDataSent(client tcp : TCPClient) {
        if let recieved = tcp.read(50){
            print("recieved: " + String(bytes: recieved, encoding: .utf8)!)
        }
    }

    //starts dispatch queue that calls itself after completion
    func tcpCycle(client tcp : TCPClient){
        DispatchQueue.main.asyncAfter(deadline: .now() + .milliseconds(50), execute: {
            self.respondToTCPDataSent(client : tcp)

            //call tcpCycle to start another dispatch Queue
            self.tcpCycle(client: tcp)
        })
    }

}

