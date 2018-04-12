//
//  LobbyPlayerManager.swift
//  soccer game
//
//  Created by rtoepfer on 3/31/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class LobbyPlayerManager{

    var players : [PlayerInfo?] = Array(repeating: nil, count: GameScene.maxPlayers)
    var scene : SKScene
    
    init(scene : SKScene){
        self.scene = scene
    }
    
    func addPlayer(playerNumber: Int, username : String ){
       
        let newPlayer = PlayerInfo(
            playerNumber : playerNumber,
            username : username,
            emoji : ChatView.defaultEmoji
        )
        
        players[playerNumber] = newPlayer
    }
    
    func playerExists(playerNumber: Int) -> Bool{
        return players[playerNumber] != nil
    }
    
    func export() -> PlayerExport {
        return PlayerExport(players: self.players)
    }
    
    func emojiChange(for player: Int, is emoji: String){
        print("player \(player) changed emoji to \(emoji)")
        players[player]?.emoji = emoji
    }
    
}

struct PlayerInfo{
    var playerNumber : Int
    var username : String
    var emoji : String
}

class PlayerExport {
    var players : [PlayerInfo?] = Array(repeating: nil, count: GameScene.maxPlayers)
    
    init(players : [PlayerInfo?]){
        self.players = players
    }
}
