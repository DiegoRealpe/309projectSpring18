//
//  MatchMakingScene.swift
//  soccer game
//
//  Created by rtoepfer on 1/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import Alamofire

class MatchMakingScene: SKScene {
    
    let totalCount = 3
    let countTime = 0.5
    
    var tcpConn:ManagedTCPConnection?
    
    var connectLabel : SKLabelNode?
    
    override func didMove(to view: SKView) {
        
        //get nodes from parent
        self.connectLabel = self.childNode(withName: "Connect Label") as? SKLabelNode
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        guard let first = touches.first else{
            return
        }
        let location = first.location(in: self)
        
        if let label = self.connectLabel, label.contains(location) {
            askServerForTCPPort()
        }
        
    }
    
    fileprivate func askServerForTCPPort(){
        self.connectLabel?.text = "connecting ..."
        
        let requestString = "http://\(CommunicationProperties.host):\(CommunicationProperties.httpport)/tcpport"
        print("sending get to \(requestString)")
        Alamofire.request(requestString, method: .get)
            .responseString(completionHandler: respondToPortHandshake(_:))
    }
    
    private func respondToPortHandshake(_ response : DataResponse<String>){
        
        if let data:Data = response.data, let str:String = String(data: data, encoding: .utf8){
            print("got response \"\(str)\"")
            
            //parse port from response string
            guard let port = Int32(str) else{
                print("did not recieve correctly formatted port sesponse, not connecting")
                return
            }
            
            self.tcpConn = ManagedTCPConnection(address : CommunicationProperties.host, port : port)
            
            let spr = SocketPacketResponder()
            //set managedTCPConnection to use spr on read
            tcpConn?.datahandler = spr.respond(data:)
            
            transitionToGameSceneWithData(spr : spr)
        }
    }
    
    func transitionToGameSceneWithData(spr : SocketPacketResponder){
        let transitionFunction = makeAddGameSceneDataFunction(spr : spr)
        self.moveToGameScene(dataFunction : transitionFunction)
    }
    
    func makeAddGameSceneDataFunction(spr : SocketPacketResponder) -> (NSMutableDictionary) -> Void{
        
        return { (dict) -> Void in
            
            dict.setValue(self.tcpConn, forKey: UserDataKeys.managedTCPConnection.rawValue)
            dict.setValue(spr, forKey: UserDataKeys.socketPacketResponder.rawValue)
            
        }
        
    }
    
    
}
