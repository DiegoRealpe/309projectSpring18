//
//  PlayerInformation.swift
//  soccer game
//
//  Created by rtoepfer on 3/1/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation


struct RawPlayerInformation{
    var playerNumber:Int
    var Name:String = "no-name"
}


class GamePlayerInformationModel {
    
    var players = [RawPlayerInformation?](repeating: nil,count: GameScene.maxPlayers)
    
    init(players : [RawPlayerInformation] = []){
        players.forEach(addPlayer(player:))
    }
    
    func addPlayer(player : RawPlayerInformation){
        //check player has valid player number
        let playerNumber = player.playerNumber
        guard (0..<GameScene.maxPlayers).contains(playerNumber) else{
            print("invalid player number. \(playerNumber) was not in range 0..<\(GameScene.maxPlayers)")
            return
        }
        if players[playerNumber] != nil {
            print("game player information model overwriting player number \(playerNumber)")
        }
        players[playerNumber] = player
    }
}
