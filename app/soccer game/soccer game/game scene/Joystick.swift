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
    
    //between -1 and 1
    var xDirection : Double
    var yDirection : Double
    
    
    init(parent : SKNode, radius : Double, startPoint : CGPoint){
        
        self.radius = radius
        self.outerCircle = makeCircleMold(radius: radius, fillColor: UIColor.blue)
        self.innerCircle = makeCircleMold(radius: radius*0.75, fillColor: UIColor.brown)
        self.parent = parent
        self.xDirection = 0
        self.yDirection = 0
        
        moveOuterTo(point: startPoint)
        
        
        parent.addChild(self.outerCircle)
        self.outerCircle.addChild(self.innerCircle)
    }
    
    func acceptNewTouch(touches: Set<UITouch>){
        let filteredTouchSet = touches.filter(isInBottomLeftQuadrant(_:))
        //right now it will make  joystick on the touch closest to the origin, this may need to change
        if let touch = closestTouchTo(touches: filteredTouchSet, node: outerCircle){
            moveOuterTo(point: touch.location(in: self.parent))
            moveInnerTo(point: CGPoint.zero)
        }
    }
    
    func isInBottomLeftQuadrant(_ touch : UITouch) -> Bool{
        let loc = touch.location(in: self.parent)
        return loc.x < 0 && loc.y < 0
    }
    
    
    
    func acceptTouchMoved(touches: Set<UITouch>){
        let touch = closestTouchTo(touches: touches, node: outerCircle)!
        
        let outerRelativeDisplayPoint = translatePointToStayInOuter(scenePoint: touch.location(in: parent))
        assignDirection()
        
        moveInnerTo(point : outerRelativeDisplayPoint)
    }
    
    func moveOuterTo(point : CGPoint){
        outerCircle.position = point
    }
    
    func moveInnerTo(point : CGPoint){
        innerCircle.position = point
    }
    
    func getDebugMessage() -> String{
        return String(xDirection) + "," + String(yDirection)
    }
    
    func translatePointToStayInOuter(scenePoint : CGPoint) -> CGPoint{
        
        let distance = Double(outerCircle.position.distanceTo(scenePoint))
        let xDiff = Double(scenePoint.x - outerCircle.position.x)
        let yDiff = Double(scenePoint.y - outerCircle.position.y)

        if distance > self.radius{
            let divider = distance/self.radius
            return CGPoint(x :xDiff/divider,y :yDiff/divider)
        }else{
            return CGPoint(x : xDiff, y : yDiff)
        }
    }
    
    //uses the position of the inner node
    func assignDirection(){
        let x = Double(self.innerCircle.position.x)
        let y = Double(self.innerCircle.position.y)
        
        self.xDirection = x/radius
        self.yDirection = y/radius
    }
    
}

//should be called only if touches is not empty
func closestTouchTo(touches : Set<UITouch>, node :SKNode) -> UITouch?{
    var iter = touches.makeIterator()
    
    if var closest = iter.next(){ //find closest if there were touches
        var closestCloseness = closest.location(in: node).distanceTo(node.position)
        
        while let next = iter.next(){
            let nextCloseness = next.location(in: node).distanceTo(node.position)
            if nextCloseness < closestCloseness{
                closest = next
                closestCloseness = nextCloseness
            }
        }
        
        return closest
    }else{
        return nil //return nil if there were no touches
    }
}


fileprivate func makeCircleMold(radius : Double, fillColor : UIColor) -> SKShapeNode{
    let circ = SKShapeNode.init(circleOfRadius : CGFloat(radius))
    circ.fillColor = fillColor
    
    return circ
}



