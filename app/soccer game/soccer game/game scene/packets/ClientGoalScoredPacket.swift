//
//  ClientGoalScoresPacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/20/18.xx
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ClientGoalScoredPacket : SendablePacket {
    
    var playerNumber : Int
    var scoringTeam : Int
    
    init(playerNum : Int, scoringTeam : Int){
        self.playerNumber = playerNum
        self.scoringTeam = scoringTeam
    }
    
    func toByteArray() -> [UInt8] {
        return [130,UInt8(playerNumber),UInt8(scoringTeam)]
    }
    
}
