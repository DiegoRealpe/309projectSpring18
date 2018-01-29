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
    var currentseconds : Int?
    var countdownAction : SKAction?
    
    
    override func didMove(to view: SKView) {
        
        self.countdownLabel = self.childNode(withName: "Time Label") as? SKLabelNode
        
        self.currentseconds = totalCount
        
//        if let cdLabel = self.countdownLabel{
//            cdLabel.text = String(totalSeconds)
//        }
        countdownAction = makeCountdownAction()
        
        startCountdown()
        
        
    }
    
    func startCountdown(){
        if let i = self.currentseconds{
            self.currentseconds = self.totalCount
            if let cdLabel = self.countdownLabel{
                cdLabel.text = String(i)
                if let act = self.countdownAction{
                    cdLabel.run(SKAction.sequence([SKAction.wait(forDuration: self.countTime),act]))
                }
            }
            
        }
        
    }
    
    func makeCountdownAction() -> SKAction {
        return SKAction.run {
            if let cdLabel = self.countdownLabel{
                if let i = self.currentseconds{
                    self.currentseconds = i - 1
                    cdLabel.text = String(i)
                    if i <= 0{
                        self.moveToGameScene()
                    }
                }
                
                if let act = self.countdownAction{
                    cdLabel.run(SKAction.sequence([SKAction.wait(forDuration: self.countTime),act]))
                }
            }
            
        }
    }

    func moveToGameScene(){
        
        if let scene = GameScene(fileNamed: "GameScene"){
            // Set the scale mode to scale to fit the window
            scene.scaleMode = .aspectFill
            
            
            // Present the scene
            if let view = self.view{
                view.presentScene(scene)
            }
        }
        
    }
    
    
}
