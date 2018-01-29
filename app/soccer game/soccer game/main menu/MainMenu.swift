//
//  MainMenu.swift
//  soccer game
//
//  Created by rtoepfer on 1/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit

class MainMenu: SKScene {
    
    var title : SKLabelNode?
    
    override func didMove(to view: SKView) {
        print("custom class running")
        
        self.title = self.childNode(withName: "TitleLabel") as? SKLabelNode
        
        fadeInLabel(label: self.title)
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        let unwrappedPractice = self.childNode(withName: "Practice Background")
        let unwrappedJoin = self.childNode(withName: "Join Background")
        
        if let t = touches.first ,let practice = unwrappedPractice, let join = unwrappedJoin{
            let point = t.location(in: self)
            if practice.contains(point){
                moveToGameScene()
                blinkLabel(label: self.title)
            }else if join.contains(point){
                moveToMatchMakingScene()
            }
        }
        
    }
    
    
    
    func fadeInLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.alpha = 0.0
            nonOptLabel.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
    func blinkLabel(label : SKLabelNode?){
        if let nonOptLabel = label{
            nonOptLabel.run(SKAction.init(named: "Pulse")!, withKey: "fadeInOut")
        }
    }
    
}
