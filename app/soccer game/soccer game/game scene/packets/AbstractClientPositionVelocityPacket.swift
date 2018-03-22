//
//  AbstractClientPositionVelocityPacket.swift
//  soccer game
//
//  Created by rtoepfer on 3/21/18.
//  Copyright © 2018 MG 6. All rights reserved.


import Foundation

class AbstractClientPositionVelocityPacket : SendablePacket {
    
    var xPosFloat = Int32(0)
    var yPosFloat = Int32(0)
    var xVelocity = Int32(0)
    var yVelocity = Int32(0)
    
    init(xPos:Int32, yPos:Int32, xV:Int32, yV:Int32)
    {
        xPosFloat = xPos
        yPosFloat = yPos
        xVelocity = xV
        yVelocity = yV
        
    }
    
    
    func toByteArray() -> [UInt8]{
        
        
        var array = [UInt8]()
        array.reserveCapacity(18)
        
        array.append(120)
        
        //xPos stuff
        let xPos = convertToUInt8(Float(xPosFloat))
        array.append(contentsOf: xPos)
        
        
        //yPos stuff
        let yPos = convertToUInt8(Float32(yPosFloat))
        array.append(contentsOf: yPos)
        
        
        //x velocity
        let xVel = convertToUInt8(Float32(xVelocity))
        array.append(contentsOf: xVel)
        
        //y velocity
        let yVel = convertToUInt8(Float32(yVelocity))
        array.append(contentsOf: yVel)
        
        return array
    }
    
}


protocol SendablePacket{
    func toByteArray() -> [UInt8]
}

