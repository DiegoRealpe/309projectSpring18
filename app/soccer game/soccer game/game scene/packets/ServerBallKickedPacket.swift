//
//  ServerBallKickedPacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/19/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ServerBallKickedPacket {
    
    var playerNumber : Int
    
    init(data : [UInt8]){
        self.playerNumber = Int(data[1])
    }
    
}
