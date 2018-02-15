//
//  MatchMakingScene.swift
//  soccer game
//
//  Created by rtoepfer on 1/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit

class MatchMakingScene: SKScene {
    
    let totalCount = 3
    let countTime = 0.5
    
    var countdownLabel : SKLabelNode?
    var currentSeconds : Int?
    var countdownAction : SKAction?
    
    override func didMove(to view: SKView) {
        
        //get nodes from parent
        self.countdownLabel = self.childNode(withName: "Time Label") as? SKLabelNode
        
        self.currentSeconds = totalCount
        
        countdownAction = makeCountdownAction()
        
        startCountdown()
        
    }
    
    func startCountdown(){
        if let i = self.currentSeconds{
            self.currentSeconds = self.totalCount
            if let cdLabel = self.countdownLabel{
                cdLabel.text = String(i)
                
                cdLabel.run(SKAction.sequence([SKAction.wait(forDuration: self.countTime),self.countdownAction!]))
                
            }
            
        }
        
    }
    
    //make countdown action makes an SKAction that accesses local fields to decide to move to the gamescene
    //or call itself and wait another interval
    func makeCountdownAction() -> SKAction {
        return SKAction.run {
            if let cdLabel = self.countdownLabel{
                
                if let i = self.currentSeconds{
                    
                    //decrenemt seconds and apply to label
                    self.currentSeconds = i - 1
                    cdLabel.text = String(i)
                    
                    //move to game scene at conclusion of countdown 
                    if i <= 0{
                        self.transitionToGameScene()
                    }
                }
                
                //restart action in a recursive manner
                if let act = self.countdownAction{
                    cdLabel.run(SKAction.sequence([SKAction.wait(forDuration: self.countTime),act]))
                }
            
            }
            
        }
    }

    
    func transitionToGameScene(){
        self.moveToGameScene(dataFunction : addGameSceneData(_:))
    }
    
    func addGameSceneData(_ dict: NSMutableDictionary){
        dict.setValue("testvalue", forKey: "test")
        dict.setValue(2, forKey: "test2")
    }
}
