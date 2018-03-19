

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
    
    static let maxPlayers = 2
    let movementSpeed = 100.0
    let packetUpdateIntervalSeconds = 0.05
    let joystickRadius = 50.0
    
    //label used for debugging, not part of final project
    private var mockPacketLabel : SKLabelNode?
    
    private var backLabel : SKLabelNode?
    private var joyStick : Joystick?
    private var ballNode : SKSpriteNode?
    private var managedTcpConnection : ManagedTCPConnection?
    
    //after didMove is called players is initialized with the exact size of maxPlayers
    private var players : [SKSpriteNode] = []
    private var playerNumber : Int?
    
    private var localPlayerStateWasUpdated = true
    
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    override func didMove(to view: SKView) {
        print("moved to game scene")
        
        configurePlayerNodes()
        self.backLabel = self.childNode(withName: "Back Label") as? SKLabelNode
        self.ballNode = self.childNode(withName: "Ball") as? SKSpriteNode
        self.mockPacketLabel = self.childNode(withName: "Mock Packet") as? SKLabelNode
        
        configureManagedTCPConnection()
        configurePacketResponder()
    }
    
    //sets players and player num to values according to the user data passed into the scene
    private func configurePlayerNodes(){
        //get player node from Players.sks
        guard let modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as? SKSpriteNode else{
            return
        }
        
        //set players to correct length with placeholders
        self.players = [SKSpriteNode](repeating : SKSpriteNode(), count: GameScene.maxPlayers)
        
        //copy model player into each index of self.players
        for i in 0..<GameScene.maxPlayers {
            players[i] = modelPlayer.copy() as! SKSpriteNode
            players[i].physicsBody = modelPlayer.physicsBody?.copy() as? SKPhysicsBody
        }
        
        //move player with the number passed into the scene into view
        if let playerNumber = self.userData?.value(forKey: UserDataKeys.playerNumber.rawValue) as? Int{
            self.playerNumber = playerNumber
            players[playerNumber].position = defaultPlayerStartingPositions[playerNumber]!
            self.addChild(players[playerNumber])
        }
        
    }
    
    //maps funtions to packet numbers to be used by the responder in a scene-specific configuration
    private func configurePacketResponder() {
        buildPacketTypeDict()
        
        if let spr = self.userData?.value(forKey: UserDataKeys.socketPacketResponder.rawValue) as? SocketPacketResponder {
            spr.packetTypeDict = self.packetTypeDict
            print(spr)
        }
    }
    
    //optionally unwrap a ManagedTcpConnection from the UserDataPassed into the scene
    private func configureManagedTCPConnection(){
        
        if let mtcp = self.userData?.value(forKey: UserDataKeys.managedTCPConnection.rawValue) as? ManagedTCPConnection {
            self.managedTcpConnection = mtcp
            
            //run update action forever
            self.run(makeUpdateAndSendSKAction())
        }
        
    }
    
    //for individual touches
    private func touchBegins(_ touch : UITouch) {
        let position = touch.location(in: self)
        
        if self.joyStick == nil && isInBottomLeftQuadrant(_ : touch) {
            self.joyStick = Joystick(parent : self, radius : self.joystickRadius, touch : touch)
        }
        
        if self.backLabel?.contains(position) == true{
            print("back to main menu")
            self.moveToScene(.mainMenu)
        }else if self.mockPacketLabel?.contains(position) == true{
        
            print("mocking command")
            let spr = SocketPacketResponder()
            spr.packetTypeDict = self.packetTypeDict
                
            let bytes : [UInt8] = [121,1,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68, 0, 0, 0, 0]
            spr.respond(data: bytes)
            
        }
    }
    
    private func setBallPositionAndVelocity(position : CGPoint, velocity : CGVector){
        guard let ball = self.ballNode else{
            print("ball was not found")
            return
        }
        
        ball.position = position
        ball.physicsBody?.velocity = velocity
    }
    
    func makeUpdateAndSendSKAction() -> SKAction {
        let packetAction = SKAction.run({
            
            if self.localPlayerStateWasUpdated {
                let packet = self.makePlayerStatePacket(playerNumber : self.playerNumber!)
                
                print("sending packet, ",packet)
                self.managedTcpConnection?.sendTCP(data: packet)
                
                self.localPlayerStateWasUpdated = false
            }
            
        })
        let waitAction = SKAction.wait(forDuration: packetUpdateIntervalSeconds)
        
        let sequenceAction = SKAction.sequence([packetAction,waitAction])
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
    
    //for individual touches
    private func touchMoved(_ touch : UITouch) {
        
    }
    
    //for individual touches
    private func touchEnded(_ touch : UITouch) {
        
        //remove joystick from scene if joystick touch
        if let js = self.joyStick, js.wasJoystickTouch(touch) {
            js.removeSelf()
            self.joyStick = nil
        }
    }
    
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        for t in touches { self.touchBegins(t) }
    }
    
    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent?) {
        
        //unwrap joystick
        if let js = self.joyStick{
            js.acceptTouchMoved(touches: touches)
            
            //capture and react to joystick position
            let dx = js.xDirection * movementSpeed
            let dy = js.yDirection * movementSpeed
            
            //set the local player to the correct velocity
            if let playerNum = self.playerNumber {
                self.players[playerNum].physicsBody!.velocity = CGVector(dx: dx, dy: dy)
                
                if let tcp = self.managedTcpConnection {
                    let packet = self.makePlayerStatePacket(playerNumber : playerNum)
                    
                    tcp.sendTCP(data: packet)
                }
                
                self.localPlayerStateWasUpdated = true
            }
        }
        
        for t in touches { self.touchMoved(t) }
    }
    
    
    
    override func touchesEnded(_ touches: Set<UITouch>, with event: UIEvent?) {
        for t in touches { self.touchEnded(t) }
    }
    
    override func touchesCancelled(_ touches: Set<UITouch>, with event: UIEvent?) {
        for t in touches { self.touchEnded(t) }
    }
    
    
    override func update(_ currentTime: TimeInterval) {
        // Called before each frame is rendered
    }
    
    
    private func buildPacketTypeDict(){
        self.packetTypeDict[121] = PacketType(dataSize: 22, handlerFunction: executePositionPacket(data:))
    }
    
    func executePositionPacket(data : [UInt8]){
        guard data.count == 22 else{
            print("executePositionPackets did not have correct data size. expected 22, was",data.count)
            return
        }
        
        print("got position packet with data:",data)
        
        let spsm = ServerPlayerStatePacket(rawData: data)
        
        let player:SKSpriteNode = selectOrAddPlayer(playerNum : spsm.playerNumber)
        
        ApplyPositionPacketToPlayer(player: player, point: spsm.position, vector: spsm.velocity)
    }
    
    //returns player node from players whith specified index
    //if the node is not a child of the SKScene it is added as a child
    private func selectOrAddPlayer(playerNum : Int) -> SKSpriteNode{
        let player:SKSpriteNode = players[Int(playerNum)]
        if player.parent != self {
            self.addChild(player)
        }
        return player
    }
    
    private func ApplyPositionPacketToPlayer(player : SKSpriteNode, point : CGPoint, vector : CGVector){
        player.position = point
        player.physicsBody?.velocity = vector
    }
    
    func isInBottomLeftQuadrant(_ touch : UITouch) -> Bool{
        let loc = touch.location(in: self)
        return loc.x < 0 && loc.y < 0
    }
}

