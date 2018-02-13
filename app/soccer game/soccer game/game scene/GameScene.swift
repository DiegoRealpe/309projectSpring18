

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
    
    //label used for debugging, not part of final project
    private var label : SKLabelNode?
    
    private var backLabel : SKLabelNode?
    private var joyStick : JoyStick?
    private var playerNode : SKSpriteNode?
    
    let movementSpeed = 100.0
    
    override func didMove(to view: SKView) { 
        
        // get optional nodes from scene
        self.label = self.childNode(withName: "//helloLabel") as? SKLabelNode
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.playerNode = self.childNode(withName:"Player Node") as? SKSpriteNode
        
        
        self.joyStick = JoyStick(parent: self, radius: 50.0, startPoint: CGPoint(x: 0, y: 0))
    }
    
    //for individual touches
    func touchDown(atPoint pos : CGPoint) {
        
        if let n = self.backLabel{
            if n.contains(pos){
                self.moveToMainMenu()
            }
        }
        
    }
    
    //for individual touches
    func touchMoved(toPoint pos : CGPoint) {
        
    }
    
    //for individual touches
    func touchUp(atPoint pos : CGPoint) {
        
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        self.joyStick?.acceptNewTouch(touches: touches)
        let str = self.joyStick!.getDebugMessage()
        
        //update label
        if let label = self.label {
            label.text = str
        }
        
        for t in touches { self.touchDown(atPoint: t.location(in: self)) }
    }
    
    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //unwrap joystick
        if let js = self.joyStick{
            js.acceptTouchMoved(touches: touches)
            let str = js.getDebugMessage()
            if let label = self.label {
                label.text = str
            }
            
            //capture and react to joystick position
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
    
}

