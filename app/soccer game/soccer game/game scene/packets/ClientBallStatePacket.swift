//
//  ClientBallStatePacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/21/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ClientBallStatePacket : AbstractClientPositionVelocityPacket {

    //combine id bytes and position/velocity data
    override func toByteArray() -> [UInt8] {
        var array = [UInt8]()
        array.append(123)
        
        let positionAndVeocityBytes = super.toByteArray()
        
        array.append(contentsOf: positionAndVeocityBytes)
        
        return array
    }
    
}

