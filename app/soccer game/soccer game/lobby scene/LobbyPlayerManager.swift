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

    var players : [PlayerInfo?] = []
    var scene : SKScene
    
    var teamPolicy : TeamPolicy
    
    init(scene : SKScene, teamPolicy : TeamPolicy){
        
        self.teamPolicy = teamPolicy
        self.players = Array(repeating: nil,count: teamPolicy.numPlayers)
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
    
    func export(playerNum : Int) -> GameScenePlayerImport.Player {
        let p = players[playerNum]!
        return GameScenePlayerImport.Player(username: p.username, playerNumber: p.playerNumber, emoji: p.emoji)
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
    var players : [PlayerInfo?]
    
    init(players : [PlayerInfo?]){
        self.players = players
    }
}
