//
//  PlayerManager.swift
//  soccer game
//
//  Created by rtoepfer on 3/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class GamePlayerManager {
    
    var teamPolicy : TeamPolicy
    
    private var players : [Player] = []
    
    var playerNumber : Int 
    
    let playerLabelRelativePosition = CGPoint(x: 0, y: 30)
    let emojiLabelRelativePosition = CGPoint(x: 0, y: -24)
    
    var modelEmojiLabel : SKLabelNode
    var modelPlayer : SKSpriteNode
    var modelUsernameLabel : SKLabelNode
    
    var lastTeam0Touched = 0
    var lastTeam1Touched = 1
    
    var scene : GameScene
    
    init(playerNumber: Int,scene : GameScene, playerImport : GameScenePlayerImport, teamPolicy : TeamPolicy){
        self.teamPolicy = teamPolicy
        self.playerNumber = playerNumber
        self.scene = scene
        self.playerNumber = playerNumber
        self.modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as! SKSpriteNode
        self.modelEmojiLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Emoji Label") as! SKLabelNode
        self.modelUsernameLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Username Label") as! SKLabelNode
        self.players = Array(repeating: Player(), count: teamPolicy.numPlayers)
        
        configurePlayerNodes(playerImport: playerImport)
    }
    
    //returns player node from players whith specified index
    //if the node is not a child of the SKScene it is added as a child
    func selectPlayer(playerNum : Int) -> SKSpriteNode{
        let player:SKSpriteNode = players[Int(playerNum)].node
        return player
    }
    
    //sets players and player num to values according to the user data passed into the scene
    fileprivate func addPlayer(importedPlayer : GameScenePlayerImport.Player) {
        let i = importedPlayer.playerNumber
        
        players[i].node = modelPlayer.copy() as! SKSpriteNode
        players[i].node.physicsBody = modelPlayer.physicsBody?.copy() as? SKPhysicsBody
        players[i].node.physicsBody?.mass = modelPlayer.physicsBody!.mass //don't know why this is necesarry
        
        print("mass of player is",players[i].node.physicsBody!.mass)
        
        scene.addChild(players[i].node)
        
        //add emoji labels
        let label = modelEmojiLabel.copy() as! SKLabelNode
        label.position = self.emojiLabelRelativePosition
        label.text = importedPlayer.emoji != nil ? importedPlayer.emoji : ""
        players[i].node.addChild(label)
        
        //add color to player
        let coloringNode = SKShapeNode(circleOfRadius: 25.0)
        let color = teamPolicy.teamColor(forPlayer: i)
        coloringNode.fillColor = color
        coloringNode.strokeColor = color
        
        self.players[i].username = importedPlayer.username
        
        players[i].node.addChild(coloringNode)
    }
    
    private func configurePlayerNodes(playerImport : GameScenePlayerImport){
        
        modelPlayer.physicsBody?.categoryBitMask = GameScene.playerCategory
        modelPlayer.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        for player in playerImport.players {
            addPlayer(importedPlayer: player)
        }
        
        addUsernameLabelsToAllPlayers(playerImport: playerImport)
        
        setToStartingPositions()
    }
    
    private func addUsernameLabelsToAllPlayers(playerImport : GameScenePlayerImport){
        for i in 0 ..< playerImport.players.count {
            let player = players[i]
            addUsernameLabelToPlayer(player.node, username: player.username)
        }
        
    }
    
    private func addUsernameLabelToPlayer(_ player : SKSpriteNode,username : String){
        let copiedLabel = self.modelUsernameLabel.copy() as! SKLabelNode
        copiedLabel.position = self.playerLabelRelativePosition
        copiedLabel.text = username
        
        player.addChild(copiedLabel)
    }
    
    
    private func addLabelToUserPlayer(){
        guard let modelLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Player Label") as? SKLabelNode else{
            return
        }
        let copiedLabel = modelLabel.copy() as! SKLabelNode
        
        let player = selectPlayer(playerNum : self.playerNumber)
        copiedLabel.position = self.playerLabelRelativePosition
        
        player.addChild(copiedLabel)
    }
    
    func removePlayerFromGame(playerNumber : Int){
        let playerToRemove = self.selectPlayer(playerNum : playerNumber)
        playerToRemove.removeFromParent()
    }

    func setToStartingPositions(){
        
        let numPlayers = scene.isLocalGame() ? 1 : GameScene.maxPlayers
        for i in 0..<numPlayers {
            players[i].node.position = teamPolicy.startingPosition(forPlayer: i)
            players[i].node.physicsBody!.velocity = .zero
        }
        
    }
    
    func recordInteractionWithBall(playerNum : Int){
        let team = self.teamPolicy.teamNumber(forPlayer: playerNum)
        if team == 0 {
            lastTeam0Touched = playerNum
        }else{
            lastTeam1Touched = playerNum
        }
    }
    
    func lastTouchForTeam(team : Int) -> Int{
        if team == 0 {
            return lastTeam0Touched
        }else{
            return lastTeam1Touched
        }
    }
    
    func usernameFor(playerNumber : Int) -> String {
        return players[playerNumber].username
    }
    
    func playerNumberFor(sprite : SKNode) -> Int?{
        for i in 0..<teamPolicy.numPlayers {
            if self.players[i].node == sprite{
                return i
            }
        }
        
        return nil
    }
    
    func countActive() -> Int {
        var active = 0
        for player in players{
            active += player.active ? 1 : 0 //increment if active
        }
        return active
    }
    
    struct Player{
        var username = "DEFAULTUSERNAME"
        var emoji = ""
        var active = true
        var node : SKSpriteNode!
    }
}
