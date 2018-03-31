//
//  ReadyManager.swift
//  soccer game
//
//  Created by rtoepfer on 3/31/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class ReadyManager{
    
    var scene : SKScene
    var readyLabel : SKLabelNode
    
    var numReady = 0
    var currentPlayerReady = false
    
    init(scene : SKScene){
        self.scene = scene
        self.readyLabel = scene.childNode(withName: "Ready Label") as! SKLabelNode
        
        updateReadyLabel()
    }
    
    func updateReadyLabel(){
        let actionText = currentPlayerReady ? "I'm not ready 😡" : "Ready Up"
        self.readyLabel.fontColor = currentPlayerReady ? UIColor.red : UIColor.green
        
        self.readyLabel.text = "\(actionText) (\(numReady)/\(GameScene.maxPlayers))"
    }
    
    func acceptTouch(touch : UITouch){
        let location = touch.location(in: scene)
        if self.readyLabel.contains(location){
            readyLabelWasTouched()
        }
    }
    
    func readyLabelWasTouched(){
        if currentPlayerReady {
            numReady -= 1
            currentPlayerReady = false
        }else{
            numReady += 1
            currentPlayerReady = true
        }
        
        updateReadyLabel()
    }
    
    func readyUserPlayer(){
        self.numReady += 1
    }
    
}
