//
//  ServerPlayerStatePacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/19/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

public class ServerPlayerStatePacket {
    
    var playerNumber : Int
    
    var position : CGPoint
    var velocity : CGVector
    
    init(rawData: [UInt8]){
        let xPosBytes = Array(rawData[2...5])
        let yPosBytes = Array(rawData[6...9])
        let xVelBytes = Array(rawData[10...13])
        let yVelBytes = Array(rawData[14...17])
        
        let xPosFloat = convertToFloat(xPosBytes).toCGFloat()
        let yPosFloat = convertToFloat(yPosBytes).toCGFloat()
        let xVelFloat = convertToFloat(xVelBytes).toCGFloat()
        let yVelFloat = convertToFloat(yVelBytes).toCGFloat()
        
        self.playerNumber = Int(rawData[1])
        self.position = CGPoint (x : xPosFloat, y: yPosFloat)
        self.velocity = CGVector(dx: xVelFloat, dy: yVelFloat)
        
    }
    
    
}
