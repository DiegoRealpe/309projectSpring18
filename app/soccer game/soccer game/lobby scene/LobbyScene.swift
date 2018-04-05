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
    
    var pm : LobbyPlayerManager?
    var rm : ReadyManager?
    var viewController: UIViewController?
    
    var chatView : ChatView!
    
    override func didMove(to view: SKView) {
        
        unpackTransitionDictionary()
        
        self.pm = LobbyPlayerManager(scene : self)
        self.rm = ReadyManager(scene: self)
        
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        
        self.pm!.addPlayer(playerNumber: playerNumber , username: "NENENEHdhe")
        
        startChatView()
    }
    
    private func startChatView(){
        DispatchQueue.main.sync {
            chatView = GameViewController.globalChatView
            print(chatView)
            chatView.loadChat()
            chatView.isHidden = false
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
    
    private func quitWasPressed() {
        print("back to main menu")
        
        
        self.moveToScene(.mainMenu)
    }
    
}
