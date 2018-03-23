//
//  ClientPlayerStatePacket.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ClientPlayerStatePacket: AbstractClientPositionVelocityPacket{
    
    
    //combine id bytes and position/velocity data
    override func toByteArray() -> [UInt8] {
        var array = [UInt8]()
        array.append(120)
        
        let positionAndVeocityBytes = super.toByteArray()
        
        array.append(contentsOf: positionAndVeocityBytes)
        
        print(array)
        
        return array
    }
    
}
