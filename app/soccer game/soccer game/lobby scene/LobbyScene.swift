//
//  LobbyScene.swift
//  soccer game
//
//  Created by rtoepfer on 3/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit
import SpriteKit
import FacebookCore
import FacebookLogin
import FBSDKLoginKit

class LobbyScene: SKScene {

    var quitLabel : SKLabelNode?
    
    var mtcp : ManagedTCPConnection!
    var spr : SocketPacketResponder!
    var playerNumber : Int!
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    var pm : LobbyPlayerManager?
    var rm : ReadyManager?
    var viewController: UIViewController?
    
    var chatView : ChatView!
    
    override func didMove(to view: SKView) {
        
        unpackTransitionDictionary()
        
        self.pm = LobbyPlayerManager(scene : self)
        self.rm = ReadyManager(scene: self)
        
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        
        self.pm!.addPlayer(playerNumber: playerNumber , username: "NENENEHdhe 🐥🇺🇸")
        
        populatePacketTypeDict()
        
        startChatView()
    }
    
    private func startChatView(){
        DispatchQueue.main.sync {
            chatView = GameViewController.globalChatView
            print(chatView)
            chatView.loadChat()
            chatView.isHidden = false
            chatView.onNewMessage = self.newLocalMessage(text:)
        }
    }
    
    private func hideChat(){
        chatView.isHidden = true
    }
    
    private func unpackTransitionDictionary(){
        self.mtcp = self.userData!.value(forKey: UserDataKeys.managedTCPConnection.rawValue) as! ManagedTCPConnection
        self.spr = self.userData!.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as! SocketPacketResponder
        self.playerNumber = self.userData!.value(forKey: UserDataKeys.playerNumber.rawValue) as! Int
    }
    
    //no multitouch support
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        guard let touch = touches.first else{
            return
        }
        
        let touchPosition = touch.location(in: self)
        if let quit = self.quitLabel, quit.contains(touchPosition){
            quitWasPressed()
        }
        
        self.rm?.acceptTouch(touch: touch)
    }
    
    private func quitWasPressed() {
        print("back to main menu")
        hideChat()
        self.moveToScene(.mainMenu)
    }
    
    private func populatePacketTypeDict(){
        self.packetTypeDict[206] = PacketType(dataSize: 82, handlerFunction: playerAddedHandler(data:))
        self.packetTypeDict[203] = PacketType(dataSize: 402, handlerFunction: chatMessageHandler(data:))
        
        self.spr.packetTypeDict = self.packetTypeDict
    }
    
    private func newLocalMessage(text : String){
        let message = OutgoingChatMessagePacket(text)
        
        print("sending size",message.toByteArray().count)
        
        self.mtcp.sendTCP(packet: message)
    }
    
    private func playerAddedHandler(data : [UInt8]){
        let player = RemotePlayerJoinedLobbyPacket(data: data)
        
        print("player added",player.playerNumber,player.username,"size is",player.username.count)
        
        self.pm?.addPlayer(playerNumber: player.playerNumber, username: player.username)
    }

    private func chatMessageHandler (data: [UInt8]){
        let message = IncommingChatMessagePacket(data: data)
        
        print("got message from",message.playerNumber,"it was",message.message)
        
        self.chatView.addRemoteMessage(message.message, from: "todo")
        
    }
    
}
