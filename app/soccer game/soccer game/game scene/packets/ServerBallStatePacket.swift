//
//  ServerBallStatePacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/21/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ServerBallStatePacket : AbstractServerPositionVelocityPacket {
    
    override init(rawData: [UInt8]){
    
        //pass position-velocity
        super.init(rawData: Array(rawData[1...16]) )
    }
    
}
