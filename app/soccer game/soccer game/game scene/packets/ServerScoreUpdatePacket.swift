//
//  ServerScoreUpdatePacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/21/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ServerScoreUpdatePacket {
    
    var team1Score : Int
    var team2Score : Int
    var scoringPlayer : Int
    
    init(raw : [UInt8]){
        team1Score = Int(raw[1])
        team2Score = Int(raw[2])
        scoringPlayer = Int(raw[3])
    }
    
}
