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
    
    private var parent : SKNode
    private var innerCircle : SKShapeNode
    private var outerCircle : SKShapeNode
    private var radius : Double
    
    
    init(parent : SKNode, radius : Double, startPoint : CGPoint){
        
        self.radius = radius
        self.outerCircle = makeCircleMold(radius: radius, fillColor: UIColor.blue)
        self.innerCircle = makeCircleMold(radius: radius*0.75, fillColor: UIColor.brown)
        self.parent = parent
        
        moveOuterTo(point: startPoint)
        
        
        parent.addChild(self.outerCircle)
        self.outerCircle.addChild(self.innerCircle)
    }
    
    func acceptNewTouch(touches: Set<UITouch>){
        //right now it will make  joystick on the touch closest to the origin, this may need to change
        let touch = closestTouchTo(touches: touches, node: outerCircle)
        
        moveOuterTo(point: touch.location(in: self.parent))
        moveInnerTo(point: CGPoint.zero)
    }
    
    
    
    func acceptTouchMoved(touches: Set<UITouch>){
        let touch = closestTouchTo(touches: touches, node: outerCircle)
        
        let outerRelativeDisplayPoint = translatePointToStayInOuter(scenePoint: touch.location(in: parent))
        
        
        moveInnerTo(point : outerRelativeDisplayPoint)
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
    
    func translatePointToStayInOuter(scenePoint : CGPoint) -> CGPoint{
        
        let distance = Double(outerCircle.position.distanceTo(scenePoint))
        let xDiff = Double(scenePoint.x - outerCircle.position.x)
        let yDiff = Double(scenePoint.y - outerCircle.position.y)
        
        if distance > radius{
            let divider = distance/self.radius
            return CGPoint(x :xDiff/divider,y :yDiff/divider)
        }else{
            return CGPoint(x : xDiff, y : yDiff)
        }
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



