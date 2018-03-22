//
//  ServerPlayerStatePacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/19/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

public class ServerPlayerStatePacket : AbstractServerPositionVelocityPacket{
    
    var playerNumber : Int
    
    override init(rawData: [UInt8]){
        
        self.playerNumber = Int(rawData[1])
        
        //pass position-velocity
        super.init(rawData: Array(rawData[2...17]) )
    }
    
}
