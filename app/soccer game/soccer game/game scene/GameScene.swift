

//
//  GameScene.swift
//  soccer game
//
//  Created by Mark Schwartz on 1/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import GameplayKit

class GameScene: SKScene {
    
    private var label : SKLabelNode?
    private var backLabel : SKLabelNode?
    private var joyStick : JoyStick?
    private var playerNode : SKSpriteNode?
    
    
    let movementSpeed = 100.0
    
    override func didMove(to view: SKView) {
        
        // Get label node from scene and store it for use later
        self.label = self.childNode(withName: "//helloLabel") as? SKLabelNode
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.playerNode = self.childNode(withName:"Player Node") as? SKSpriteNode
        
        fadeInMainLabel()
        
        self.joyStick = JoyStick(parent: self, radius: 50.0, startPoint: CGPoint(x: 0, y: 0))
    }
    
    func touchDown(atPoint pos : CGPoint) {
        
        if let n = self.backLabel{
            if n.contains(pos){
                self.moveToMainMenu()
            }
        }
        
    }
    
    func touchMoved(toPoint pos : CGPoint) {
        
    }
    
    func touchUp(atPoint pos : CGPoint) {
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //self.joyStick?.moveOuterTo(touches : touches)
        self.joyStick?.acceptNewTouch(touches: touches)
        let str = self.joyStick!.getDebugMessage()
        
        if let label = self.label {
            label.text = str
            label.run(SKAction.init(named: "Pulse")!, withKey: "fadeInOut")
        }
        
        for t in touches { self.touchDown(atPoint: t.location(in: self)) }
    }
    
    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        if let js = self.joyStick{
            js.acceptTouchMoved(touches: touches)
            let str = js.getDebugMessage()
            if let label = self.label {
                label.text = str
            }
            let dx = js.xDirection * movementSpeed
            let dy = js.yDirection * movementSpeed
            self.playerNode?.physicsBody?.velocity = CGVector(dx: dx, dy: dy)
            
        }
        
        
        for t in touches { self.touchMoved(toPoint: t.location(in: self)) }
    }
    
    override func touchesEnded(_ touches: Set<UITouch>, with event: UIEvent?) {
        for t in touches { self.touchUp(atPoint: t.location(in: self)) }
    }
    
    override func touchesCancelled(_ touches: Set<UITouch>, with event: UIEvent?) {
        for t in touches { self.touchUp(atPoint: t.location(in: self)) }
    }
    
    
    override func update(_ currentTime: TimeInterval) {
        // Called before each frame is rendered
    }
    
    private func fadeInMainLabel() {
        if let label = self.label {
            label.alpha = 0.0
            label.run(SKAction.fadeIn(withDuration: 2.0))
        }
    }
    
}

