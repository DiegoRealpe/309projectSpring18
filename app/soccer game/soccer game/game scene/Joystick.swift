//
//  Joystick.swift
//  soccer game
//
//  Created by rtoepfer on 1/31/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import SpriteKit

class Joystick{
    
    private var parent : SKNode
    private var innerCircle : SKShapeNode
    private var outerCircle : SKShapeNode
    private var radius : Double
    private var joystickTouch : UITouch
    
    //between -1 and 1 such that x^2 + y^2 <= 1
    //exposed to be read as outputs
    var xDirection : Double
    var yDirection : Double
    
    init(parent : SKNode, radius : Double, touch : UITouch){
        
        self.radius = radius
        self.outerCircle = makeCircle(radius: radius, fillColor: UIColor.blue)
        self.innerCircle = makeCircle(radius: radius*0.75, fillColor: UIColor.brown)
        self.parent = parent
        self.xDirection = 0
        self.yDirection = 0
        self.joystickTouch = touch
        
        //move joystick to initial position
        moveOuterTo(point: touch.location(in: parent))
        positionNewJoystick(touch: touch)
        
        //add joystick nodes to their parents
        parent.addChild(self.outerCircle)
        self.outerCircle.addChild(self.innerCircle)
        
    }
    
    private func positionNewJoystick(touch: UITouch){
        
        moveOuterTo(point: touch.location(in: self.parent))
        moveInnerTo(point: CGPoint.zero)
    }
    
    //uses the fact that touchesBegan and touchesMoved ..ect all use the same literal object
    func wasJoystickTouch(_ touch : UITouch) -> Bool {
        return touch == self.joystickTouch
    }
    
    func removeSelf() {
        self.outerCircle.removeFromParent()
    }
    
    func acceptTouchMoved(touches: Set<UITouch>){
        
        //respond to touch closest to the center of the joysick
        if let touch = touches.first(where: wasJoystickTouch(_:)){
            
            //translate point to keep position in
            let finalPoint = translatePointToStayInOuter(scenePoint: touch.location(in: parent))
            
            moveInnerTo(point : finalPoint)
            
            setXYDirection()
        }
    }
    
    private func moveOuterTo(point : CGPoint){
        outerCircle.position = point
    }
    
    private func moveInnerTo(point : CGPoint){
        innerCircle.position = point
    }
    
    func getDebugMessage() -> String{
        return String(xDirection) + "," + String(yDirection)
    }
    
    //returns position for inner node so that x^2 + y^2 <= radius of outer
    private func translatePointToStayInOuter(scenePoint : CGPoint) -> CGPoint{
        
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
    private func setXYDirection(){
        let x = Double(self.innerCircle.position.x)
        let y = Double(self.innerCircle.position.y)
        
        //keep output independent of joystick size by dividing out raius
        self.xDirection = x/radius
        self.yDirection = y/radius
    }
    
    private func closestTouchTo(touches : Set<UITouch>, node :SKNode) -> UITouch?{
        return touches.min(by: self.areTouchesInAscendingOrderByDistanceToCenter(first:second:))
    }
    
    //used to filter set for min in accordance to swift's built-in set functionality
    private func areTouchesInAscendingOrderByDistanceToCenter(first : UITouch, second : UITouch) -> Bool{
        let compareNode = self.outerCircle
        
        //compare to (0,0) since locations are centered relative to the compareNode
        let firstDistance = first.location(in: compareNode).distanceTo(.zero)
        let secondDistance = second.location(in: compareNode).distanceTo(.zero)
        
        return firstDistance < secondDistance
    }
}


fileprivate func makeCircle(radius : Double, fillColor : UIColor) -> SKShapeNode{
    let circ = SKShapeNode.init(circleOfRadius : CGFloat(radius))
    circ.fillColor = fillColor
    
    return circ
}



