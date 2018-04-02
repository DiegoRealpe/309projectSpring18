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
    private var httpHandshakeLabel:SKLabelNode?
    private var sendTCPLabel:SKLabelNode?
    private var stopTCPLabel:SKLabelNode?
    private var sendHelloLabel:SKLabelNode?
    
    private var tcpConn : ManagedTCPConnection?
    
    //sandbox here
    override func didMove(to view: SKView) {
        
        print("at comm test screen")
        
        self.httpHandshakeLabel = self.childNode(withName: "Http Handshake") as? SKLabelNode
        self.sendTCPLabel = self.childNode(withName: "Send Tcp") as? SKLabelNode
        self.stopTCPLabel = self.childNode(withName: "Close Tcp") as? SKLabelNode
        self.sendHelloLabel = self.childNode(withName: "Send Hello") as? SKLabelNode
        
        //testHttp() //use refer string response to recievedResponse
    }
    
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        if let touch = touches.first{
            if let tcpLabel = self.sendTCPLabel, tcpLabel.contains(touch.location(in: self)){
                //self.tcpConn = ManagedTCPConnection(address : "localhost", port : 7234)
            }
            else if let stopLabel = self.stopTCPLabel, stopLabel.contains(touch.location(in: self)){
                self.tcpConn?.stop()
            }
            else if let sendLabel = self.sendHelloLabel, sendLabel.contains(touch.location(in: self)){
                self.tcpConn?.sendTCP(message: "hello")
            }
            else if let httpLabel = self.httpHandshakeLabel, httpLabel.contains(touch.location(in: self)){
                askServerForTCPPort()
            }
        }
    }
    
    
    private func testHttp(){
        Alamofire.request("http://localhost:8080", method: .get)
                .responseString(completionHandler: recievedResponse(_:))
    }
    
    //to be passed to alamofire to handle string response
    //needs future multithreading support
    func recievedResponse(_ response : DataResponse<String>){
        
        if let data:Data = response.data, let str:String = String(data: data, encoding: .utf8){
            print("got response " + str)
            
            self.httpHandshakeLabel!.text = str
        }
    }
    
    fileprivate func askServerForTCPPort(){
        let requestString = "http://\(CommunicationProperties.gameServiceHost):\(CommunicationProperties.gameServiceHttpPort)/tcpport"
        print("sending get to \(requestString)")
        Alamofire.request(requestString, method: .get)
            .responseString(completionHandler: respondToPortHandshake(_:))
    }
    
    private func respondToPortHandshake(_ response : DataResponse<String>){
        
        if let data:Data = response.data, let str:String = String(data: data, encoding: .utf8){
            print("got response \"\(str)\"")
            
            guard let port = Int32(str) else{
                print("did not recieve correctly formatted port sesponse, not connecting")
                return
            }
            
            self.httpHandshakeLabel!.text = str
            self.tcpConn = ManagedTCPConnection(address : CommunicationProperties.gameServiceHost, port : port)
        }
    }
    
    
}

