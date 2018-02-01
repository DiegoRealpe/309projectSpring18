//
//  Joystick.swift
//  soccer game
//
//  Created by rtoepfer on 1/31/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class JoyStick{
    
    
    private var innerCircle : SKShapeNode
    private var outerCircle : SKShapeNode
    private var startPoint : CGPoint
    
    init(parent : SKScene, radius : Double, startPoint : CGPoint){
        
        self.startPoint = startPoint
        self.outerCircle = makeCircleMold(radius: radius, fillColor: UIColor.blue)
        self.innerCircle = makeCircleMold(radius: radius*0.75, fillColor: UIColor.brown)
        
        moveOuterTo(point: startPoint)
        
        
        parent.addChild(self.outerCircle)
        self.outerCircle.addChild(self.innerCircle)
    }
    
    func acceptNewTouch(touches: Set<UITouch>, parent : SKNode){
        //right now it will make  joystick on the touch closest to the origin, this may need to change
        let touch = closestTouchTo(touches: touches, node: outerCircle)
        
        moveOuterTo(point: touch.location(in: parent))
        moveInnerTo(point: CGPoint.zero)
    }
    
    
    func acceptTouchMoved(touches: Set<UITouch>){
        let touch = closestTouchTo(touches: touches, node: outerCircle)
        
        print(touch.location(in: outerCircle))
        
        moveInnerTo(point: touch.location(in: outerCircle))
    }
    
    func moveOuterTo(point : CGPoint){
        outerCircle.position = point
    }
    
    func moveInnerTo(point : CGPoint){
        innerCircle.position = point
    }
    
    func getDebugMessage() -> String{
        //return String(describing : self.outerCircle.position.x)+","+String(describing : self.outerCircle.position.x)
        return String(describing : self.outerCircle.position)
    }
    
}

//should be called only if touches is not empty
func closestTouchTo(touches : Set<UITouch>, node :SKNode) -> UITouch{
    var iter = touches.makeIterator()
    
    var closest = iter.next()!
    var closestCloseness = closest.location(in: node).distanceTo(node.position)
    
    while let next = iter.next(){
        let nextCloseness = next.location(in: node).distanceTo(node.position)
        if nextCloseness < closestCloseness{
            closest = next
            closestCloseness = nextCloseness
        }
    }
    
    return closest
}

fileprivate func makeCircleMold(radius : Double, fillColor : UIColor) -> SKShapeNode{
    let circ = SKShapeNode.init(circleOfRadius : CGFloat(radius))
    circ.fillColor = fillColor
    
    return circ
}



