

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
    
    let maxPlayers = 2
    let movementSpeed = 100.0
    let offScreen = CGPoint(x :10000, y :10000)
    let packetUpdateIntervalSeconds = 0.05
    
    //label used for debugging, not part of final project
    private var mockPacketLabel : SKLabelNode?
    
    private var backLabel : SKLabelNode?
    private var joyStick : JoyStick?
    private var ballNode : SKSpriteNode?
    private var managedTcpConnection : ManagedTCPConnection?
    
    private var players : [SKSpriteNode] = []
    private var playerNumber : Int?
    
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    override func didMove(to view: SKView) {
        
        // get optional nodes from scene
        configurePlayerNodes()
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.ballNode = self.childNode(withName: "Ball") as? SKSpriteNode
        self.joyStick = JoyStick(parent: self, radius: 50.0, startPoint: CGPoint(x: 0, y: 0))
        self.mockPacketLabel = self.childNode(withName: "Mock Packet") as? SKLabelNode
        
        print(self.mockPacketLabel!)
        
        configureManagedTCPConnection()
        configurePacketResponder()
    }
    
    func configurePlayerNodes(){
        
        guard let modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as? SKSpriteNode else{
            return
        }
        
        //init players with placeholders
        self.players = [SKSpriteNode](repeating : SKSpriteNode(), count: maxPlayers)
        
        for i in 0..<maxPlayers {
            players[i] = modelPlayer.copy() as! SKSpriteNode
            players[i].physicsBody = modelPlayer.physicsBody?.copy() as? SKPhysicsBody
        }
        
        print(players[1])
        
        if let playerNumber = self.userData?.value(forKey: UserDataKeys.playerNumber.rawValue) as? Int{
            players[playerNumber].position = CGPoint(x : 100, y : -100)
            self.playerNumber = playerNumber
            self.addChild(players[playerNumber])
        }
        
    }
    
    func configurePacketResponder() {
        buildPacketTypeDict()
        
        if let spr = self.userData?.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as? SocketPacketResponder {
            spr.packetTypeDict = self.packetTypeDict
            print(spr)
        }
    }
    
    func configureManagedTCPConnection(){
        if let mtcp = self.userData?.value(forKey: UserDataKeys.managedTCPConnection.rawValue) as? ManagedTCPConnection {
            self.managedTcpConnection = mtcp
            
            //run update action forever
            self.run(makeUpdateAndSendSKAction())
        }
    }
    
    //for individual touches
    func touchDown(atPoint pos : CGPoint) {
        
        if self.backLabel?.contains(pos) == true{
            print("back to main menu")
            self.moveToMainMenu()
        }else if self.mockPacketLabel?.contains(pos) == true{
        
            print("mocking command")
            let spr = SocketPacketResponder()
            spr.packetTypeDict = self.packetTypeDict
                
            let bytes : [UInt8] = [121,1,0,0,0,0,0,0,0,0,152, 78, 154, 68,152, 78, 154, 68]
            spr.respond(data: bytes)
            
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
            //startIntervalSendLoop()
            if let playerNum = self.playerNumber {
                self.players[playerNum].physicsBody!.velocity = CGVector(dx: dx, dy: dy)
                
                if let tcp = self.managedTcpConnection {
                    let packet = self.makePlayerStatePacket(playerNumber : playerNum)
                    
                    tcp.sendTCP(data: packet)
                }
            }
        }
        
        
        
        for t in touches { self.touchMoved(toPoint: t.location(in: self)) }
    }
    
    func makeUpdateAndSendSKAction() -> SKAction {
        let readerAction = SKAction.run({
            let packet = self.makePlayerStatePacket(playerNumber : self.playerNumber!)
            
            self.managedTcpConnection?.sendTCP(data: packet)
        })
        let waitAction = SKAction.wait(forDuration: packetUpdateIntervalSeconds)
        let sequenceAction = SKAction.sequence([readerAction,waitAction])
        
        return SKAction.repeatForever(sequenceAction)
    }
    
    
    func makePlayerStatePacket(playerNumber : Int)-> [UInt8]
    {
        let chosenPlayer = self.players[playerNumber]
        let position = chosenPlayer.position
        let velocity = chosenPlayer.physicsBody?.velocity
        
        
        let playerPacket = ClientPlayerStatePacket(xPos: Int32(position.x), yPos: Int32(position.y), xV: Int32(velocity!.dx), yV: Int32(velocity!.dy))
        
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
        
        let player:SKSpriteNode = selectOrAddPlayer(playerNum : Int(playerNum))
        
        ApplyPositionPacketToPlayer(player: player, point: position, vector: velocity)
    }
    
    func selectOrAddPlayer(playerNum : Int) -> SKSpriteNode{
        let player:SKSpriteNode = players[Int(playerNum)]
        if player.parent != self {
            self.addChild(player)
        }
        return player
    }
    
    func ApplyPositionPacketToPlayer(player : SKSpriteNode, point : CGPoint, vector : CGVector){
        player.position = point
        player.physicsBody?.velocity = vector
    }
}

