//
//  ClientPlayerStatePacket.swift
//  soccer game
//
//  Created by Mark Schwartz on 2/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
class ClientPlayerStatePacket
{
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
        array.append(120)
        
        var xPos = UInt8(xPosFloat & 0x1000)
        array.append(xPos)
        
         xPos = UInt8(xPosFloat & 0x0100)
        array.append(xPos)
        
         xPos = UInt8(xPosFloat & 0x0010)
        array.append(xPos)
        
         xPos = UInt8(xPosFloat & 0x0001)
        array.append(xPos)
        
        
        //yPos stuff
        var yPos = UInt8(yPosFloat & 0x1000)
        array.append(yPos)
        
        yPos = UInt8(yPosFloat & 0x0100)
        array.append(yPos)
        
        yPos = UInt8(yPosFloat & 0x0010)
        array.append(yPos)
        
        yPos = UInt8(yPosFloat & 0x0001)
        array.append(yPos)
        
        //x velocity
        var xVel = UInt8(xVelocity & 0x1000)
        array.append(xVel)
        
        xVel = UInt8(xVelocity & 0x0100)
        array.append(xVel)
        
        xVel = UInt8(xVelocity & 0x0010)
        array.append(xVel)
        
        xVel = UInt8(xVelocity & 0x0001)
        array.append(xVel)
        
        //y velocity
        var yVel = UInt8(yVelocity & 0x1000)
        array.append(yVel)
        
        yVel = UInt8(yVelocity & 0x0100)
        array.append(yVel)
        
        yVel = UInt8(yVelocity & 0x0010)
        array.append(yVel)
        
        yVel = UInt8(yVelocity & 0x0001)
        array.append(yVel)
        
        return array
    }
    
    
}
