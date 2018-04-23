//
//  GameNotifier.swift
//  soccer game
//
//  Created by rtoepfer on 4/23/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class GameNotifier {
    
    private var scene : SKScene
    private var displayLabel : SKLabelNode
    
    init(scene: SKScene){
        self.scene = scene
        self.displayLabel = scene.childNode(withName: "Notification Label") as! SKLabelNode
    }
    
    func displayMessage(_ str: String){
        
        let action = makeShowAction(text: str)
        displayLabel.run(action)
        
    }
    
    
    private func makeShowAction(text: String) -> SKAction{
        displayLabel.isHidden = false
        displayLabel.text = text
        let act0 = SKAction.fadeIn(withDuration: 0.2)
        let act1 = SKAction.wait(forDuration: 0.6)
        let act2 = SKAction.fadeOut(withDuration: 0.2)
        
        return SKAction.sequence([act0,act1,act2])
        
    }
    
    
}
