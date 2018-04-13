//
//  GameScenePlayerImport.swift
//  soccer game
//
//  Created by rtoepfer on 4/13/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

struct GameScenePlayerImport {
    
    var players : [Player]
    
    struct Player {
        var username : String
        var playerNumber : Int
        var emoji : String?
    }
}
