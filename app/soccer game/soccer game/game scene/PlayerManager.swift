//
//  PlayerManager.swift
//  soccer game
//
//  Created by rtoepfer on 3/29/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class PlayerManager {
    
    var players : [SKSpriteNode] = []
    var playerNumber : Int
    let playerLabelRelativePosition = CGPoint(x: 0, y: 30)
    
    var scene : SKScene
    
    init(playerNumber: Int,scene : SKScene){
        self.playerNumber = playerNumber
        self.scene = scene
        
        configurePlayerNodes()
    }
    
    //returns player node from players whith specified index
    //if the node is not a child of the SKScene it is added as a child
    func selectPlayer(playerNum : Int) -> SKSpriteNode{
        let player:SKSpriteNode = players[Int(playerNum)]
        return player
    }
    
    //sets players and player num to values according to the user data passed into the scene
    private func configurePlayerNodes(){
        //get player node from Players.sks
        guard let modelPlayer = SKScene(fileNamed : "Players")?.childNode(withName : "Player Node") as? SKSpriteNode else{
            return
        }
        modelPlayer.physicsBody?.categoryBitMask = GameScene.playerCategory
        
        modelPlayer.physicsBody?.contactTestBitMask = GameScene.ballCategory
        
        
        
        //set players to correct length with placeholders
        self.players = [SKSpriteNode](repeating : SKSpriteNode(), count: GameScene.maxPlayers)
        
        //copy model player into each index of self.players
        for i in 0..<GameScene.maxPlayers {
            players[i] = modelPlayer.copy() as! SKSpriteNode
            players[i].physicsBody = modelPlayer.physicsBody?.copy() as? SKPhysicsBody
            
            players[i].position = defaultPlayerStartingPositions[i]!
            scene.addChild(players[i])
        }
        
        addLabelToUserPlayer()
    }
    
    func removePlayerFromGame(playerNumber : Int){
        let playerToRemove = self.selectPlayer(playerNum : playerNumber)
        playerToRemove.removeFromParent()
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
}
