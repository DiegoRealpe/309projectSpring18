

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
    let forceUpdateWaits = 1
    
    var waitsSinceLastPlayerUpdate = 0
    var waitsSinceLastBallUpdate = 0
   
    var kickCoolDownTime : Double = 0.0
    let kickCoolDownInterval = 0.5
    
    //after didMove is called players is initialized with the exact size of maxPlayers
    var pm : GamePlayerManager!
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
        configureCollisions()
        configurePlayerManager()
        
        //give all children of the north bounds(all the bounds) the same physics category
        for child in (northBound?.children)!
        {
            child.physicsBody?.categoryBitMask = GameScene.boundsCategory
            child.physicsBody?.contactTestBitMask = GameScene.ballCategory
        }
        
        //scoreboard stuff
        self.redTeamScore = self.childNode(withName: "Left Team Score") as? SKLabelNode
        self.blueTeamScore = self.childNode(withName: "Right Team Score") as? SKLabelNode
        scoreBoard = ScoreBoard(redTeamLabel: redTeamScore!, blueTeamLabel: blueTeamScore!)
        
        
        configureManagedTCPConnection()
        configurePacketResponder()
        
        kickButton = KickButton(scene: self)
    }
    
    func configurePlayerManager(){
        let playerImport = self.userData!.value(forKey: UserDataKeys.gameSecnePlayerImport.rawValue) as! GameScenePlayerImport
        let playerNumber = self.userData!.value(forKey: UserDataKeys.playerNumber.rawValue) as! Int
        self.pm = GamePlayerManager(playerNumber: playerNumber, scene: self, playerImport: playerImport)
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
        self.northBound = self.childNode(withName: "North Bound") as? SKSpriteNode
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
        else if isScorekeeper() {
            if(otherCategory == GameScene.leftGoalCategory)
            {
                scoreBoard?.redTeamScored()
                print("Left Goal Scored")
                let scorePacket = ClientGoalScoredPacket(playerNum: 0, scoringTeam: 0)
                managedTcpConnection?.sendTCP(packet: scorePacket)
            }
            else if(otherCategory == GameScene.rightGoalCategory)
            {
                scoreBoard?.blueTeamScored()
                print("Right Goal Scored")
                let scorePacket = ClientGoalScoredPacket(playerNum: 0, scoringTeam: 0)
                managedTcpConnection?.sendTCP(packet: scorePacket)
            }
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
        }else{
            print("🐝🐝🐝 game is in practice mode as no ManagedTCP connection was passed into the scene")
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
        
        if self.joyStick == nil && (isInBottomLeftQuadrant(_ : touch) || true) {
            self.joyStick = Joystick(parent : self, radius : self.joystickRadius, touch : touch)
        }
        if self.kickButton.contains(position){
            doLocalKick()
        }
        if self.quitLabel?.contains(position) == true{
            self.quitWasPressed()
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
                    print("sending",packet)
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
        self.packetTypeDict[204] = PacketType(dataSize: 2, handlerFunction: stub(_:))
        self.packetTypeDict[127] = PacketType(dataSize: 1, handlerFunction: executeHostAssignmentPacket(data:))
        self.packetTypeDict[131] = PacketType(dataSize: 4, handlerFunction: executeScoreUpdatePacket(data:))
        self.packetTypeDict[134] = PacketType(dataSize: 2, handlerFunction: executeKickPacket(data:))
    }
    
    func executePlayerPositionPacket(data : [UInt8]){
        guard data.count == 22 else{
            print("executePlayerPositionPackets did not have correct data size. expected 22, was",data.count)
            return
        }
        
        print("got player position packet with data:",data)
        
        let spsp = ServerPlayerStatePacket(rawData: data)
        
        print("player: \(spsp.playerNumber), i am \(pm.playerNumber) x: \(spsp.position.x) y: \(spsp.position.y)")
        
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
        
        self.pm.removePlayerFromGame(playerNumber: Int(data[1]))
    }
    
    func executeHostAssignmentPacket(data : [UInt8]){
        self.isHost = true
        print("got host assignment packet 😎")
    }
    
    func doLocalKick(){
        doKick(playerNum: self.pm.playerNumber)
    }
    
    func doKick(playerNum : Int){
        
        let playerPosition = self.pm!.selectPlayer(playerNum: playerNum).position
        let distanceBetweenBallAndPlayer : Float = self.ballNode!.position.distanceTo(playerPosition)
        let now : Double = CACurrentMediaTime();
        
        print("distnce was: \(distanceBetweenBallAndPlayer))")
        
        print("time is \(now) and coolDownTime is \(kickCoolDownTime)")
        if distanceBetweenBallAndPlayer < GameScene.maxKickDistance, now > self.kickCoolDownTime {
            let vector : CGVector =  playerPosition.vectorTo(self.ballNode!.position,ofMagnitude: 300)
            
            print("kick",vector)
            self.ballNode!.physicsBody!.velocity = vector
            sendKickPacket()
            
            self.kickCoolDownTime = now + kickCoolDownInterval
        }
    }
    
    private func executeKickPacket(data : [UInt8]){
        let packet = ServerBallKickedPacket(data: data)
        
        doKick(playerNum: packet.playerNumber)
    }
    
    private func sendKickPacket(){
        guard let mtcp = self.managedTcpConnection else {
            return //the game is local
        }
        
        let packet = ClientBallKickedPacket()
        
        mtcp.sendTCP(packet: packet)
    }
    
    func isScorekeeper() -> Bool {
        return self.isHost || isPractice()
    }
    
    func isPractice() -> Bool {
        return self.managedTcpConnection == nil
    }
    
    private func stub(_: [UInt8]) {
        return
    }
    
    func executeScoreUpdatePacket(data : [UInt8]){
        let packet = ServerScoreUpdatePacket(raw : data)
        
        self.scoreBoard?.forceScore(team1: packet.team1Score, team2: packet.team2Score)
        self.pm.setInitialPositions()
        self.setBallToStartingPosition()
    }
    
    func setBallToStartingPosition(){
        self.ballNode?.position = .zero
    }
}

