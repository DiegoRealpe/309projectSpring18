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
    
    //for simplicity this scene currently only explicitly supplrts 1 touch
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        guard let first = touches.first else{
            return
        }
        let location = first.location(in: self)
        
        if let label = self.connectLabel, label.contains(location) {
            askServerForTCPPort()
        }
    }
    
    //start handshake where server gives TCP port to use for sockets
    private func askServerForTCPPort(){
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
            
            let spr = SocketPacketResponder()
            spr.packetTypeDict = makePacketTypeDict(spr:spr)
            
            //set ManagedTCPConnection to use spr as responder
            self.tcpConn = ManagedTCPConnection(address : CommunicationProperties.host, port : port, dataHandler : spr.respond(data:))
        }
    }
    
    //for interfacing with the SocketPacketResponder
    private func makePacketTypeDict(spr : SocketPacketResponder) -> [UInt8:PacketType]{
        return [
            122: PacketType(dataSize: 2, handlerFunction: { (data) in
                self.recievePlayerNumberCode(data: data, spr: spr)
            })
        ]
    }
    
    
    private func recievePlayerNumberCode(data : [UInt8],spr : SocketPacketResponder){
        guard data.count == 2 else {
            print("did not recieve correct player code size, expected 2, was",data.count)
            return
        }
        
        print("assigned player number :",data[1])

        transitionToGameSceneWithData(spr : spr,playerNum: data[1])
    }
    
    private func transitionToGameSceneWithData(spr : SocketPacketResponder, playerNum : UInt8){
        let transitionFunction = makeAddGameSceneDataFunction(spr : spr, playerNum : playerNum)
        self.moveToScene(.gameScene, dataFunction : transitionFunction)
    }
    
    private func makeAddGameSceneDataFunction(spr : SocketPacketResponder, playerNum : UInt8) -> (NSMutableDictionary) -> Void{
        
        return { (dict) -> Void in
            
            dict.setValue(playerNum, forKey: UserDataKeys.playerNumber.rawValue)
            dict.setValue(self.tcpConn, forKey: UserDataKeys.managedTCPConnection.rawValue)
            dict.setValue(spr, forKey: UserDataKeys.socketPacketResponder.rawValue)
            
        }
        
    }
    
    
}
