

//
//  GameScene.swift
//  soccer game
//
//  Created by Mark Schwartz on 1/28/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import SpriteKit
import GameplayKit

class GameScene: SKScene , SKPhysicsContactDelegate {
    
    static let maxPlayers = 2
    static let maxKickDistance : Float = 80.0
    
    let movementSpeed = 100.0
    let packetUpdateIntervalSeconds = 0.05
    let joystickRadius = 50.0
    
    //label used for debugging, not part of final project
    var mockPacketLabel : SKLabelNode?
    
    var quitLabel : SKLabelNode?
    var joyStick : Joystick?
    var kickButton : KickButton!
    var ballNode : SKSpriteNode?
    var managedTcpConnection : ManagedTCPConnection?
    var leftGoal: SKSpriteNode?
    var rightGoal: SKSpriteNode?
    var northBound : SKSpriteNode?
    var redTeamScore: SKLabelNode?
    var blueTeamScore:SKLabelNode?
    var scoreBoard:ScoreBoard?
    let forceUpdateWaits = 50
    var waitsSinceLastPlayerUpdate = 0
    var waitsSinceLastBallUpdate = 0
   
    
    //after didMove is called players is initialized with the exact size of maxPlayers
    var pm : PlayerManager?
    var isHost = false
    
    var localPlayerStateWasUpdated = true
    var localBallStateWasUpdates = true //make true when contact is detected
    
    var packetTypeDict : [UInt8:PacketType] = [:]
    
    static let boundsCategory:UInt32 = 0b1
    static let playerCategory:UInt32 = 0b1 << 1
    static let ballCategory:UInt32 = 0b1 << 2;
    static let leftGoalCategory:UInt32 = 0b1 << 3;
    static let rightGoalCategory:UInt32 = 0b1 << 4;
    
    override func didMove(to view: SKView) {
        print("moved to game scene")
        
        getNodesFromScene()
        
        let playerNumber = self.lookupPlayerNumber()
        self.northBound = self.childNode(withName: "North Bound") as? SKSpriteNode
        configureCollisions()
        
        self.pm = PlayerManager(playerNumber : playerNumber,scene : self)
        self.isHost = self.pm!.playerNumber == 0 //todo make more complex logic
        
        //give all children of the north bounds(all the bounds) the same physics category
        for child in (northBound?.children)!
        {
            child.physicsBody?.categoryBitMask = GameScene.boundsCategory
            child.physicsBody?.contactTestBitMask = GameScene.ballCategory
        }
        
        self.mockPacketLabel = self.childNode(withName: "Mock Packet") as? SKLabelNode
        
        //scoreboard stuff
        self.redTeamScore = self.childNode(withName: "Left Team Score") as? SKLabelNode
        self.blueTeamScore = self.childNode(withName: "Right Team Score") as? SKLabelNode
        scoreBoard = ScoreBoard(redTeamLabel: redTeamScore!, blueTeamLabel: blueTeamScore!)
        
        
        configureManagedTCPConnection()
        configurePacketResponder()
        
        kickButton = KickButton(scene: self)
    }
    
    fileprivate func configureCollisions() {
        self.physicsWorld.contactDelegate = self
        
        self.ballNode?.physicsBody?.categoryBitMask = GameScene.ballCategory
        self.ballNode?.physicsBody?.contactTestBitMask = GameScene.playerCategory | GameScene.leftGoalCategory |  GameScene.rightGoalCategory
        
        self.leftGoal?.physicsBody?.categoryBitMask = GameScene.leftGoalCategory
        self.leftGoal?.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        self.rightGoal?.physicsBody?.categoryBitMask = GameScene.rightGoalCategory
        self.rightGoal?.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        self.northBound?.physicsBody?.categoryBitMask = GameScene.boundsCategory
        self.northBound?.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
    }
    
    fileprivate func getNodesFromScene() {
        self.quitLabel = self.childNode(withName: "Quit Label") as? SKLabelNode
        self.ballNode = self.childNode(withName: "Ball") as? SKSpriteNode
        self.leftGoal = self.childNode(withName: "Left Goal") as? SKSpriteNode
        self.rightGoal = self.childNode(withName: "Right Goal") as? SKSpriteNode
    }
    
    
    func didBegin(_ contact: SKPhysicsContact) {
        let firstCategory:UInt32 = contact.bodyA.categoryBitMask//know what category this object is in
        let secondCategory:UInt32 = contact.bodyB.categoryBitMask
        
        if(firstCategory == GameScene.ballCategory || secondCategory == GameScene.ballCategory)//we know one of the objects is the ball
        {
            //get node that isnt the ball
            let otherNode:SKNode = (firstCategory == GameScene.ballCategory) ? contact.bodyB.node! : contact.bodyA.node!
            ballDidCollide(with: otherNode)
        }
        
    }
    func ballDidCollide(with other:SKNode)
    {
        let otherCategory = other.physicsBody?.categoryBitMask
        if (otherCategory == GameScene.boundsCategory)//could add more later to see if ball/player collide
        {
            print("Ball Hit Bounds")
            localBallStateWasUpdates = true
        }
        else if(otherCategory == GameScene.playerCategory)
        {
            print("Player hit ball")
            localBallStateWasUpdates = true
        }
        else if(otherCategory == GameScene.leftGoalCategory)
        {
            scoreBoard?.redTeamScored()
            print("Left Goal Scored")
        }
        else if(otherCategory == GameScene.rightGoalCategory)
        {
            scoreBoard?.blueTeamScored()
            print("Right Goal Scored")
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
    fileprivate func quitWasPressed() {
        print("back to main menu")
        
        //if tcp connection exists, disconnect from server
        if let mtcp = self.managedTcpConnection {
            mtcp.sendTCP(data: [125]) //send packet to disconnect
            mtcp.stop()
        }
        
        self.moveToScene(.mainMenu)
    }
    
    private func touchBegins(_ touch : UITouch) {
        let position = touch.location(in: self)
        
        if self.joyStick == nil && isInBottomLeftQuadrant(_ : touch) {
            self.joyStick = Joystick(parent : self, radius : self.joystickRadius, touch : touch)
        }
        if self.kickButton.contains(position){
            doKick()
        }
        if self.quitLabel?.contains(position) == true{
            self.quitWasPressed()
        }
        if self.mockPacketLabel?.contains(position) == true{
        
            print("mocking command")
            let spr = SocketPacketResponder()
            spr.packetTypeDict = self.packetTypeDict
                
            let bytes : [UInt8] = [124,0,0,0,0,0,0,0,0,152, 78, 154, 68, 152, 78, 154, 68,0,0,0,0]
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
    
    fileprivate func sendBallStatePacketIfNecesarry() {
        guard self.isHost else{
            return
        }
        
        if (self.localBallStateWasUpdates && self.ballNode != nil) || self.waitsSinceLastBallUpdate >= self.forceUpdateWaits {
            let packet = self.makeBallStatePacket()
            
            print("sending ball packet, ",packet.toByteArray())
            self.managedTcpConnection?.sendTCP(packet: packet)
            
            self.localBallStateWasUpdates = false
            self.waitsSinceLastBallUpdate = 0
        }
        else{
            self.waitsSinceLastBallUpdate += 1
        }
    }
    
    fileprivate func sendPlayerStatePacketIfNecesarry() {
        if self.localPlayerStateWasUpdated || self.waitsSinceLastPlayerUpdate >= self.forceUpdateWaits  {
            let packet = self.makePlayerStatePacket(playerNumber : self.pm!.playerNumber)
            
            print("sending player packet, ",packet.toByteArray())
            self.managedTcpConnection?.sendTCP(packet: packet)
            
            self.localPlayerStateWasUpdated = false
            
            self.waitsSinceLastPlayerUpdate = 0
        }else{
            self.waitsSinceLastPlayerUpdate += 1
        }
    }
    
    func makeUpdateAndSendSKAction() -> SKAction {
        let packetAction = SKAction.run({
            self.sendPlayerStatePacketIfNecesarry()
            self.sendBallStatePacketIfNecesarry()
        })
        
        //run the action
        let waitAction = SKAction.wait(forDuration: packetUpdateIntervalSeconds)
        let sequenceAction = SKAction.sequence([packetAction,waitAction])
        return SKAction.repeatForever(sequenceAction)
    }
    
    func makePlayerStatePacket(playerNumber : Int)-> ClientPlayerStatePacket
    {
        let chosenPlayer = self.pm!.selectPlayer(playerNum : playerNumber)
        let position = chosenPlayer.position
        let velocity = chosenPlayer.physicsBody?.velocity
        
        
        let playerPacket = ClientPlayerStatePacket(xPos: Int32(position.x), yPos: Int32(position.y), xV: Int32(velocity!.dx), yV: Int32(velocity!.dy))
        
        return playerPacket
    }
    
    func makeBallStatePacket() -> ClientBallStatePacket {
        
        let position = ballNode!.position
        let velocity = ballNode!.physicsBody!.velocity
        
        print(Int32(position.x))
        print(Int32(position.y))
        print(Int32(velocity.dx))
        print(Int32(velocity.dy))
        
        return ClientBallStatePacket(xPos: Int32(position.x), yPos: Int32(position.y), xV: Int32(velocity.dx), yV: Int32(velocity.dy))
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
            if let playerNum = self.pm?.playerNumber {
                self.pm?.selectPlayer(playerNum: playerNum).physicsBody!.velocity = CGVector(dx: dx, dy: dy)
                
                if let tcp = self.managedTcpConnection {
                    let packet = self.makePlayerStatePacket(playerNumber : playerNum)
                    tcp.sendTCP(packet: packet)
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
        self.packetTypeDict[121] = PacketType(dataSize: 22, handlerFunction: executePlayerPositionPacket(data:))
        self.packetTypeDict[124] = PacketType(dataSize: 21, handlerFunction: executeBallPositionPacket(data:))
        self.packetTypeDict[126] = PacketType(dataSize: 2, handlerFunction: executePlayerLeftGamePacket(data:))
    }
    
    func executePlayerPositionPacket(data : [UInt8]){
        guard data.count == 22 else{
            print("executePlayerPositionPackets did not have correct data size. expected 22, was",data.count)
            return
        }
        
        print("got player position packet with data:",data)
        
        let spsp = ServerPlayerStatePacket(rawData: data)
        
        let player:SKSpriteNode = self.pm!.selectPlayer(playerNum : spsp.playerNumber)
        
        ApplyPositionPacketToPlayer(player: player, point: spsp.position, vector: spsp.velocity)
    }
    
    func executeBallPositionPacket(data : [UInt8]){
        guard data.count == 21 else{
            print("executeBallPositionPackets did not have correct data size. expected 21, was",data.count)
            return
        }
        
        print("got ball position packet with data:",data)
        
        let sbsp = ServerBallStatePacket(rawData: data)
        
        self.ballNode?.position = sbsp.position
        self.ballNode?.physicsBody?.velocity = sbsp.velocity
    }
    
    
    private func ApplyPositionPacketToPlayer(player : SKSpriteNode, point : CGPoint, vector : CGVector){
        player.position = point
        player.physicsBody?.velocity = vector
    }
    
    func isInBottomLeftQuadrant(_ touch : UITouch) -> Bool{
        let loc = touch.location(in: self)
        return loc.x < 0 && loc.y < 0
    }
    
    func executePlayerLeftGamePacket(data : [UInt8]){
        guard data.count == 2 else{
            print("executePlayerLeftGamePacket did not have correct data size. expected 2, was",data.count)
            return
        }
        print("got player left game packet with data:",data)
        
        self.pm!.removePlayerFromGame(playerNumber: Int(data[1]))
    }
    
    func lookupPlayerNumber() -> Int {
        return self.userData!.value(forKey: UserDataKeys.playerNumber.rawValue) as! Int
    }
    
    func doKick(){
        let playerPosition = self.pm!.selectPlayer(playerNum: self.pm!.playerNumber).position
        
        let distanceBetweenBallAndPlayer : Float = self.ballNode!.position.distanceTo(playerPosition)
        
        print("distnce was: \(distanceBetweenBallAndPlayer))")
        
        if distanceBetweenBallAndPlayer < GameScene.maxKickDistance {
            let vector : CGVector =  playerPosition.vectorTo(self.ballNode!.position,ofMagnitude: 300)
            
            print("kick",vector)
            self.ballNode!.physicsBody!.velocity = vector
        }
    }
}

