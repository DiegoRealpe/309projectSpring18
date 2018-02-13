//
//  CommTestScreen.swift
//  soccer game
//
//  Created by rtoepfer on 2/9/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import Alamofire
import SwiftSocket

//for a screen to test network communications in an enviroment with no reprecussions, will not be in final producrt
class CommTestScreen: SKScene {
    
    //label for output
    private var debugLabel:SKLabelNode?
    private var sendTCPLabel:SKLabelNode?
    
    //private var client:TCPClient?
    
    //sandbox here
    override func didMove(to view: SKView) {
        
        print("at comm test screen")
        
        self.debugLabel = self.childNode(withName: "Debug Label") as? SKLabelNode
        self.sendTCPLabel = self.childNode(withName: "Send Tcp") as? SKLabelNode
        
        
        
        //testHttp() //use refer string response to recievedResponse
    }
    
    private func startTCPClient(){
        //self.client = TCPClient(address: "localhost", port: 7234)
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        if let touch = touches.first, let tcpLabel = self.sendTCPLabel{
            if(tcpLabel.contains(touch.location(in: self))){
                sendTcp()
            }
        }
    }
    
    func sendTcp(){
        let client = TCPClient(address: "localhost", port: 7234)
        client.connect(timeout: 2).logError()
        //client.send(string: "hello server\n").logError()

        
        DispatchQueue.main.asyncAfter(deadline: .now() + .seconds(1), execute: {
            if let recieved = client.read(50){
                print(String(bytes: recieved, encoding: .utf8)!)
                client.send(string: "ryan ").logError()
                self.step2(client: client)
            }
        })
        
    }
    
    func step2(client : TCPClient){
        
        DispatchQueue.main.asyncAfter(deadline: .now() + .seconds(1), execute: {
            if let recieved = client.read(100){
                print(String(bytes: recieved, encoding: .utf8)!)
                
                //client.send(string : ":quit").logError()
                client.close()
            }
        })
        
        
    }
    
    func tcpLogic(client : TCPClient){
        
    }
    
    fileprivate func testHttp(){
        Alamofire.request("http://localhost:8080", method: .get)
                .responseString(completionHandler: recievedResponse(_:))
    }
    
    //to be passed to alamofire to handle string response
    //needs future multithreading support
    func recievedResponse(_ response : DataResponse<String>){
        
        if let data:Data = response.data, let str:String = String(data: data, encoding: .utf8){
            print("got response " + str)
            
            self.debugLabel!.text = str
        }
    }
    
    
}

