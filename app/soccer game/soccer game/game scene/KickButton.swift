//
//  KickButton.swift
//  soccer game
//
//  Created by rtoepfer on 4/2/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class KickButton {
    
    var scene : SKScene
    var kickButton : SKShapeNode =  SKShapeNode(circleOfRadius: 30)
    
    init(scene : SKScene){
        self.scene = scene
        addButtonNode()
    }
 
    func addButtonNode(){
        kickButton.alpha = 0.6
        kickButton.fillColor = UIColor.blue
        kickButton.strokeColor = UIColor.blue
        kickButton.position = CGPoint(x: 170, y: -100)
        kickButton.zPosition = 1
        
        let label = SKLabelNode(text: "Kick")
        label.position = CGPoint(x: 0, y: -10)
        kickButton.addChild(label)
        
        scene.addChild(kickButton)
    }
    
    func contains(_ position : CGPoint) -> Bool{
        return kickButton.contains(position)
    }
}


