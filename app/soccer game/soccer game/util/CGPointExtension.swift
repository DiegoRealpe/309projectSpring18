
//
//  CGPointExtensions.swift
//  soccer game
//
//  Created by rtoepfer on 2/1/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

extension CGPoint{
    
    func distanceTo(_ other : CGPoint) -> Float {
        
        let xDiff = Float(other.x - self.x)
        let yDiff = Float(other.y - self.y)
        
        return sqrt(xDiff*xDiff + yDiff*yDiff)
    }
    
    func vectorTo(_ other : CGPoint, multipliedBy : Double) -> CGVector{
        let xDiff = Double(other.x - self.x)
        let yDiff = Double(other.y - self.y)
        
        return CGVector(dx: xDiff * multipliedBy, dy: yDiff * multipliedBy)
    }
    
    func vectorTo(_ other : CGPoint, ofMagnitude : Double) -> CGVector{
        let xDiff = Double(other.x - self.x)
        let yDiff = Double(other.y - self.y)
        
        let naturalMagnitude = sqrt(xDiff*xDiff + yDiff*yDiff)
        
        return CGVector(dx: xDiff * ofMagnitude / naturalMagnitude, dy: yDiff * ofMagnitude / naturalMagnitude)
    }
    
}

