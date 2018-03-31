//
//  LobbyScene.swift
//  soccer game
//
//  Created by rtoepfer on 3/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import UIKit
import SpriteKit

class LobbyScene: SKScene {

    var quitLabel : SKLabelNode?
    
    var pm : LobbyPlayerManager?
    var rm : ReadyManager?
    
    override func didMove(to view: SKView) {
        self.pm = LobbyPlayerManager(scene : self)
        self.rm = ReadyManager(scene: self)
        
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        
        self.pm!.addPlayer(playerNumber: 0, username: "NENENEHdhe")
        self.pm!.addPlayer(playerNumber: 1, username: "POIHURHRHUR")
        
        self.pm!.removePlayer(playerNumber: 0)
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
