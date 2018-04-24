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
    
    var pm : LobbyPlayerManager!
    var rm : ReadyManager!
    var viewController: UIViewController?
    
    var chatView : ChatView!
    
    override func didMove(to view: SKView) {
        
        unpackTransitionDictionary()
        
        let tp = unpackTeamPolicy()
        self.pm = LobbyPlayerManager(scene : self, teamPolicy: tp)
        
        self.rm = ReadyManager(scene: self, onReady: self.localPlayerReadied, onUnready: self.localPlayerUnreadied, lpm : self.pm)
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        
        self.pm!.addPlayer(playerNumber: playerNumber , username: localUsername())
        
        populatePacketTypeDict()
        
        startChatView()
    }
    
    private func startChatView(){
        DispatchQueue.main.sync {
            chatView = GameViewController.globalChatView
            print(chatView)
            chatView.lpm = self.pm
            chatView.loadChat()
            chatView.isHidden = false
            chatView.textInput.delegate = chatView
            chatView.player0Emoji.delegate = chatView
            chatView.player1Emoji.delegate = chatView
            chatView.player2Emoji.delegate = chatView
            chatView.player3Emoji.delegate = chatView
            chatView.onNewMessage = self.newLocalMessage(text:)
            chatView.onEmojiChange = self.onLocalEmojiChange(for:is:)
            
            chatView.addPlayer(playerNum: self.playerNumber ,username: localUsername() ,emojiEditable: true)
        }
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
    
    private func handle207(data: [UInt8]){
        print("quitting player was",data[1])
        
        DispatchQueue.main.sync(execute: self.chatView.hideAndClose)
        self.mtcp.stop()
        self.moveToScene(.mainMenu)
    }
    
    private func quitWasPressed() {
        print("back to main menu")
        self.mtcp.sendTCP(data : [125])
        self.mtcp.stop()
        
        self.chatView.hideAndClose()
        self.moveToScene(.mainMenu)
    }
    
    private func populatePacketTypeDict(){
        self.packetTypeDict[206] = PacketType(dataSize: 82, handlerFunction: playerAddedHandler(data:))
        self.packetTypeDict[203] = PacketType(dataSize: 402, handlerFunction: chatMessageHandler(data:))
        self.packetTypeDict[204] = PacketType(dataSize: 2, handlerFunction: remotePlayerReadiedHandler(data:))
        self.packetTypeDict[205] = PacketType(dataSize: 2, handlerFunction: remotePlayerUnreadiedHandler(data:))
        self.packetTypeDict[207] = PacketType(dataSize: 2, handlerFunction: handle207(data:))
        self.packetTypeDict[209] = PacketType(dataSize: 26, handlerFunction: remoteEmojiChanged(data:))
        self.packetTypeDict[122] = PacketType(dataSize: 1, handlerFunction: movingToGameHandler(data:))
        
        self.spr.packetTypeDict = self.packetTypeDict
    }
    
    private func newLocalMessage(text : String){
        let message = OutgoingChatMessagePacket(text)
        
        print("sending size",message.toByteArray().count)
        
        self.mtcp.sendTCP(packet: message)
    }
    
    func localPlayerReadied(){
        print("local player readied")
        self.mtcp.sendTCP(data: [200])
    }
    
    func localPlayerUnreadied(){
        print("local player unreadied")
        self.mtcp.sendTCP(data: [201])
    }
    
    func remotePlayerReadiedHandler(data : [UInt8]){
        print("remote player readied",data)
        self.rm?.readyRemote(num: Int(data[1]))
    }
    
    func remotePlayerUnreadiedHandler(data : [UInt8]){
        print("remote player unreadied",data)
        self.rm?.unreadyRemote(num: Int(data[1]))
    }
    
    private func playerAddedHandler(data : [UInt8]){
        let player = RemotePlayerJoinedLobbyPacket(data: data)
        
        print("player added",player.playerNumber,player.username,"size is",player.username.count)
        
        self.pm?.addPlayer(playerNumber: player.playerNumber, username: player.username)
        
        DispatchQueue.main.sync {
            self.chatView.addPlayer(playerNum: player.playerNumber, username: player.username,emojiEditable: false)
        }
        
    }

    private func chatMessageHandler(data: [UInt8]){
        let message = IncommingChatMessagePacket(data: data)
        
        print("got message from",message.playerNumber,"it was",message.message)
        
        self.chatView.addRemoteMessage(message.message, from: "todo")
        
    }
    
    func remoteEmojiChanged(data: [UInt8]){
        let emojiChange = IncomingEmojiChangedPacket(data)
        
        self.pm!.emojiChange(for: emojiChange.playerNumber, is: emojiChange.emoji)
        
        DispatchQueue.main.async {
            self.chatView.changeEmoji(playerNumber: emojiChange.playerNumber, emoji: emojiChange.emoji)
        }
    }
    
    func onLocalEmojiChange(for player: Int, is emoji: String){
        let emojiPacket = OutgoingEmojiChangedPacket(emoji: emoji)
        self.mtcp.sendTCP(packet: emojiPacket)
        
        self.pm!.emojiChange(for: player, is: emoji)
    }
    
    func movingToGameHandler(data : [UInt8]){
        print("moving to game scene")
        
        DispatchQueue.main.sync(execute: self.chatView.hideAndClose )
        
        self.moveToScene(.gameScene, dataFunction: makeGameSceneTransitionDisctionary(dict:))
    }
    
    func makeGameSceneTransitionDisctionary(dict: NSMutableDictionary){
        var playerImport = GameScenePlayerImport(players: [])
        for i in 0..<2 {
            let player = self.pm.export(playerNum: i)
            playerImport.players.append(player)
        }
        
        dict.setValue(playerImport, forKey: UserDataKeys.gameSecnePlayerImport.rawValue)
        dict.setValue(self.playerNumber, forKey: UserDataKeys.playerNumber.rawValue)
        dict.setValue(self.mtcp, forKey: UserDataKeys.managedTCPConnection.rawValue)
        dict.setValue(self.spr, forKey: UserDataKeys.socketPacketResponder.rawValue)
        dict.setValue(TwoPlayerTeamPolicy(), forKey: UserDataKeys.teamPolicy.rawValue)
        
    }
    
    func localUsername() -> String{
        if let name = LocalPlayerInfo.username {
            return name
        }else{
            return "Default Username"
        }
    }
    
    private func unpackTeamPolicy() -> TeamPolicy {
        return self.userData!.value(forKey: UserDataKeys.teamPolicy.rawValue) as! TeamPolicy
    }
}
