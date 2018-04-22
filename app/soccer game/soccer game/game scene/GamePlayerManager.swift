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
    
    var players : [SKSpriteNode] = []
    var playerNumber : Int
    
    let playerLabelRelativePosition = CGPoint(x: 0, y: 30)
    let emojiLabelRelativePosition = CGPoint(x: 0, y: -24)
    
    var modelEmojiLabel : SKLabelNode
    var modelPlayer : SKSpriteNode
    var modelUsernameLabel : SKLabelNode
    
    var scene : GameScene
    
    init(playerNumber: Int,scene : GameScene, playerImport : GameScenePlayerImport){
        self.playerNumber = playerNumber
        self.scene = scene
        self.playerNumber = playerNumber
        self.modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as! SKSpriteNode
        self.modelEmojiLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Emoji Label") as! SKLabelNode
        self.modelUsernameLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Username Label") as! SKLabelNode
        self.players = Array(repeating: SKSpriteNode(), count: 2)
        
        configurePlayerNodes(playerImport: playerImport)
    }
    
    //returns player node from players whith specified index
    //if the node is not a child of the SKScene it is added as a child
    func selectPlayer(playerNum : Int) -> SKSpriteNode{
        let player:SKSpriteNode = players[Int(playerNum)]
        return player
    }
    
    //sets players and player num to values according to the user data passed into the scene
    fileprivate func addPlayer(importedPlayer : GameScenePlayerImport.Player) {
        let i = importedPlayer.playerNumber
        
        players[i] = modelPlayer.copy() as! SKSpriteNode
        players[i].physicsBody = modelPlayer.physicsBody?.copy() as? SKPhysicsBody
        players[i].physicsBody?.mass = modelPlayer.physicsBody!.mass //don't know why this is necesarry
        
        print("mass of player is",players[i].physicsBody!.mass)
        
        scene.addChild(players[i])
        
        //add username and emoji labels
        let label = modelEmojiLabel.copy() as! SKLabelNode
        label.position = self.emojiLabelRelativePosition
        label.text = importedPlayer.emoji != nil ? importedPlayer.emoji : ""
        players[i].addChild(label)
        
        //add color to player
        let coloringNode = SKShapeNode(circleOfRadius: 25.0)
        let team = teamNumber(forPlayer: i)
        if team == 0{
            coloringNode.fillColor = #colorLiteral(red: 0.9254902005, green: 0.2352941185, blue: 0.1019607857, alpha: 1)
            coloringNode.strokeColor = #colorLiteral(red: 0.9254902005, green: 0.2352941185, blue: 0.1019607857, alpha: 1)
        }else {
            coloringNode.fillColor = #colorLiteral(red: 0.01680417731, green: 0.1983509958, blue: 1, alpha: 1)
            coloringNode.strokeColor = #colorLiteral(red: 0.01680417731, green: 0.1983509958, blue: 1, alpha: 1)
        }
        players[i].addChild(coloringNode)
    }
    
    private func configurePlayerNodes(playerImport : GameScenePlayerImport){
        
        modelPlayer.physicsBody?.categoryBitMask = GameScene.playerCategory
        modelPlayer.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        //set players to correct length with placeholders
        self.players = [SKSpriteNode](repeating : SKSpriteNode(), count: GameScene.maxPlayers)
        
        for player in playerImport.players {
            addPlayer(importedPlayer: player)
        }
        
        addUsernameLabelsToAllPlayers(playerImport: playerImport)
        
        setToStartingPositions()
    }
    
    private func addUsernameLabelsToAllPlayers(playerImport : GameScenePlayerImport){
        for i in 0 ..< playerImport.players.count {
            let player = players[i]
            let username = playerImport.players[i].username
            addUsernameLabelToPlayer(player, username: username)
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
        
        let numPlayers = scene.isPractice() ? 1 : GameScene.maxPlayers
        for i in 0..<numPlayers {
            players[i].position = defaultPlayerStartingPositions[i]!
            players[i].physicsBody!.velocity = .zero
        }
        
    }
    
    func teamNumber(forPlayer num: Int) -> Int{
        return num
    }
}
