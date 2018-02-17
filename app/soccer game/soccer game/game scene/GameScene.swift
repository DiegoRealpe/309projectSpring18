

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
    private var ballNode : SKSpriteNode?
    
    let movementSpeed = 100.0
    
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    override func didMove(to view: SKView) { 
        
        // get optional nodes from scene
        self.label = self.childNode(withName: "Hello Label") as? SKLabelNode
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.playerNode = self.childNode(withName:"Player Node") as? SKSpriteNode
        self.ballNode = self.childNode(withName: "Ball") as? SKSpriteNode
        self.joyStick = JoyStick(parent: self, radius: 50.0, startPoint: CGPoint(x: 0, y: 0))
        
        self.buildPacketTypeDict()
        self.userData?.value(forKey: )
        
        
    }
    
    //for individual touches
    func touchDown(atPoint pos : CGPoint) {
        
        if let n = self.backLabel{
            if n.contains(pos){
                self.moveToMainMenu()
                
            }
        }
    }
    
    func setBallPositionAndVelocity(position : CGPoint, velocity : CGVector){
        guard let ball = self.ballNode else{
            print("ball was not found")
            return
        }
        
        ball.position = position
        ball.physicsBody?.velocity = velocity
        
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
           
            
            let positionTuple = self.playerNode?.position
            let playerState = ClientPlayerStatePacket.init(xPos:Int32(positionTuple!.x) , yPos: Int32(positionTuple!.y), xV: Int32(dx), yV: Int32(dy))
            var playerByteArray = playerState.toByteArray()
            
            var tcpConn : ManagedTCPConnection?
            
            tcpConn = ManagedTCPConnection(address : "proj-309-mg-6.cs.iastate.edu", port : 5543)
            tcpConn?.sendTCP(message: "OK")
        
            
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
    
    
    func buildPacketTypeDict(){
        self.packetTypeDict[121] = PacketType(dataSize: 17, handlerFunction: executePositionPacket(data:))
    }
    
    func executePositionPacket(data : [UInt8]){
        guard data.count == 17 else{
            print("executePositionPackets did not have correct data size. expected 17, was",data.count)
            return
        }
        
        let playerNum = data[1]
        let xPosBytes = Array(data[2...5])
        let yPosBytes = Array(data[6...9])
        let xVelBytes = Array(data[10...13])
        let yVelBytes = Array(data[14...17])
        
        
        let xPosFloat = convertToFloat(xPosBytes).toCGFloat()
        let yPosFloat = convertToFloat(yPosBytes).toCGFloat()
        let xVelFloat = convertToFloat(xVelBytes).toCGFloat()
        let yVelFloat = convertToFloat(yVelBytes).toCGFloat()
        
        let position = CGPoint (x : xPosFloat, y: yPosFloat)
        let velocity = CGVector(dx: xVelFloat, dy: yVelFloat)
        
        let player:SKSpriteNode = selectPlayer(num: playerNum)
        
        ApplyPositionPacketToPlayer(player: player, point: position, vector: velocity)
    }
    
    func selectPlayer(num : UInt8) -> SKSpriteNode{
        //todo: add checking for valid player number
        return self.playerNode!
    }
    
    
    func ApplyPositionPacketToPlayer(player : SKSpriteNode, point : CGPoint, vector : CGVector){
        player.position = point
        player.physicsBody?.velocity = vector
    }
}

