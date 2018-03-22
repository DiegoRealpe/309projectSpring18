//
//  AbstractServerPositionVelocitPacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/21/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

import SpriteKit

public class AbstractServerPositionVelocityPacket {
    
    var position : CGPoint
    var velocity : CGVector
    
    init(rawData: [UInt8]){
        let xPosBytes = Array(rawData[0...3])
        let yPosBytes = Array(rawData[4...7])
        let xVelBytes = Array(rawData[8...11])
        let yVelBytes = Array(rawData[12...15])
        
        let xPosFloat = convertToFloat(xPosBytes).toCGFloat()
        let yPosFloat = convertToFloat(yPosBytes).toCGFloat()
        let xVelFloat = convertToFloat(xVelBytes).toCGFloat()
        let yVelFloat = convertToFloat(yVelBytes).toCGFloat()
        
        self.position = CGPoint (x : xPosFloat, y: yPosFloat)
        self.velocity = CGVector(dx: xVelFloat, dy: yVelFloat)
        
    }
    
    
}
