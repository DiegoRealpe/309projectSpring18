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
    
    var pm : LobbyPlayerManager?
    var rm : ReadyManager?
    var viewController: UIViewController?
    
    override func didMove(to view: SKView) {
        
        self.pm = LobbyPlayerManager(scene : self)
        self.rm = ReadyManager(scene: self)
        
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        
        self.pm!.addPlayer(playerNumber: 0, username: "NENENEHdhe")
        self.pm!.addPlayer(playerNumber: 1, username: "POIHURHRHUR")
        
        self.pm!.removePlayer(playerNumber: 0)
        
        presentChatSubview()
    }
    
    
    func presentChatSubview(){
        DispatchQueue.main.async {
//            let chatView =
//            chatView.center = self.view!.center
//            chatView.frame = CGRect(x: 100, y: 100, width: 100, height: 100)
//
//            self.view!.insertSubview(chatView, at: 0)

            //self.view!.addSubview(chatView)
            //self.view!.bringSubview(toFront: chatView)
            
            var txtField: UITextField = UITextField(frame: CGRect(x: 0, y: 0, width: 200, height: 30));
            txtField.backgroundColor = UIColor(red: 125/255, green: 125/250, blue: 125/250, alpha: 0.5)
            txtField.text = "hello"

            self.view!.addSubview(txtField)
            self.view!.superview!.insertSubview(txtField, aboveSubview: self.view!)

            print(self.view!.subviews)
            print(txtField.description)
            print(txtField.superview)

        }
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
