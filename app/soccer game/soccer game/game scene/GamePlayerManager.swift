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
    
    var scene : GameScene
    
    init(playerNumber: Int,scene : GameScene, playerImport : GameScenePlayerImport){
        self.playerNumber = playerNumber
        self.scene = scene
        self.playerNumber = playerNumber
        self.modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as! SKSpriteNode
        self.modelEmojiLabel = SKScene(fileNamed : "Players")?.childNode(withName : "Emoji Label") as! SKLabelNode
        
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
        
        players[i].position = defaultPlayerStartingPositions[i]!
        
        print("mass of player is",players[i].physicsBody!.mass)
        
        scene.addChild(players[i])
        
        let label = modelEmojiLabel.copy() as! SKLabelNode
        label.position = self.emojiLabelRelativePosition
        label.text = importedPlayer.emoji != nil ? importedPlayer.emoji : ""
        players[i].addChild(label)
    }
    
    private func configurePlayerNodes(playerImport : GameScenePlayerImport){
        
        modelPlayer.physicsBody?.categoryBitMask = GameScene.playerCategory
        modelPlayer.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        //set players to correct length with placeholders
        self.players = [SKSpriteNode](repeating : SKSpriteNode(), count: GameScene.maxPlayers)
        
        for player in playerImport.players {
            self.addPlayer(importedPlayer: player)
        }
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

}
