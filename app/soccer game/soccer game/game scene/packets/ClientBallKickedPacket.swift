//
//  ClientBallKickedPacket.swift
//  soccer game
//
//  Created by rtoepfer on 4/19/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

class ClientBallKickedPacket : SendablePacket {
    
    func toByteArray() -> [UInt8] {
        return [133]
    }
    
}
