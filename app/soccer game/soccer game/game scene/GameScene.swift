

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
    private var mockPacketLabel : SKLabelNode?
    
    private var backLabel : SKLabelNode?
    private var joyStick : JoyStick?
    private var playerNode : SKSpriteNode?
    private var ballNode : SKSpriteNode?
    private var managedTcpConnection : ManagedTCPConnection?
    
    let movementSpeed = 100.0
    
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    
    override func didMove(to view: SKView) {
        
        // get optional nodes from scene
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.playerNode = self.childNode(withName:"Player Node") as? SKSpriteNode
        self.ballNode = self.childNode(withName: "Ball") as? SKSpriteNode
        self.joyStick = JoyStick(parent: self, radius: 50.0, startPoint: CGPoint(x: 0, y: 0))
        self.mockPacketLabel = self.childNode(withName: "Mock Packet") as? SKLabelNode
        
        print(self.mockPacketLabel!)
        
        configureManagedTCPConnection()
        configurePacketResponder()
        
    }
    
    func configurePacketResponder() {
        buildPacketTypeDict()
        
        if let spr = self.userData?.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as? SocketPacketResponder {
            spr.packetTypeDict = self.packetTypeDict
            print(spr)
        }
    }
    
    func configureManagedTCPConnection(){
        if let mtcp = self.userData?.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as? ManagedTCPConnection {
            self.managedTcpConnection = mtcp
            print(mtcp)
        }
    }
    
    //for individual touches
    func touchDown(atPoint pos : CGPoint) {
        
        if self.backLabel?.contains(pos) == true{
            print("back to main menu")
            self.moveToMainMenu()
        }else if self.mockPacketLabel?.contains(pos) == true{
            print("touched the mock")
            if let spr = self.userData?.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as? SocketPacketResponder {
                
                print("mocking command")
                
                let bytes : [UInt8] = [121,0,0,0,0,0,0,0,0,0,88,78,67,33,99,23,123,45]
                spr.respond(data: bytes)
                
                
            }else{
                print("did not find socket packet responder")
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
        
        for t in touches { self.touchDown(atPoint: t.location(in: self)) }
    }
    
    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //unwrap joystick
        if let js = self.joyStick{
            js.acceptTouchMoved(touches: touches)
            
            //capture and react to joystick position
            let dx = js.xDirection * movementSpeed
            let dy = js.yDirection * movementSpeed
            self.playerNode?.physicsBody?.velocity = CGVector(dx: dx, dy: dy)
            
            if let tcp = self.managedTcpConnection {
                let packet = self.makePlayerStatePacket()
                
                tcp.sendTCP(data: packet)
            }
        }
        
        
        
        for t in touches { self.touchMoved(toPoint: t.location(in: self)) }
    }
    
    func sendStateToServer(tcp : ManagedTCPConnection){
        let bytes = makePlayerStatePacket()
    }
    
    
    func makePlayerStatePacket()-> [UInt8]
    {
        let posTuple = self.playerNode?.position
        let velTuple = self.playerNode?.physicsBody?.velocity
        
        
        let playerPacket = ClientPlayerStatePacket(xPos: Int32(posTuple!.x), yPos: Int32(posTuple!.y), xV: Int32(velTuple!.dx), yV: Int32(velTuple!.dy))
        
        return playerPacket.toByteArray()
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
        self.packetTypeDict[121] = PacketType(dataSize: 18, handlerFunction: executePositionPacket(data:))
    }
    
    func executePositionPacket(data : [UInt8]){
        guard data.count == 18 else{
            print("executePositionPackets did not have correct data size. expected 17, was",data.count)
            return
        }
        
        print(data)
        
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

