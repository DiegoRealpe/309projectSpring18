//
//  TeamPolicies.swift
//  soccer game
//
//  Created by rtoepfer on 4/22/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation
import UIKit

/*
 Team policies dictate the configuration related aspects of multiplayer and single-
 player games. Such properties include, team colors, which teams players belong to and
 other similar items.
 */
protocol TeamPolicy {
    var numPlayers : Int {get}
    
    func teamNumber(forPlayer num: Int) -> Int
    func teamColor(forPlayer num: Int) -> UIColor
    func teamColor(forTeam team: Int) -> UIColor
    func startingPosition(forPlayer num : Int) -> CGPoint
}

class PracticeTeamPolicy : TeamPolicy {
    
    private let color = #colorLiteral(red: 0.1987203089, green: 1, blue: 0.2482636255, alpha: 1)
    
    let numPlayers = 1
    
    func teamNumber(forPlayer num: Int) -> Int {
        return 0
    }
    
    func teamColor(forPlayer num: Int) -> UIColor {
        return color
    }
    
    func teamColor(forTeam team: Int) -> UIColor {
        return color
    }
    
    func startingPosition(forPlayer num: Int) -> CGPoint {
        return CGPoint(x: -100,y: -100)
    }
}

class TwoPlayerTeamPolicy : TeamPolicy {
    private let color0 = #colorLiteral(red: 0.9254902005, green: 0.2352941185, blue: 0.1019607857, alpha: 1)
    private let color1 = #colorLiteral(red: 0.01680417731, green: 0.1983509958, blue: 1, alpha: 1)
    private let startingPositions = [ //I recognize that an array might be more efficient
        0 : CGPoint(x: -100,y: 100), //but a dictionary better expresses the let's purpose
        1 : CGPoint(x: 100,y: -100),
    ]
    
    let numPlayers = 2
    
    func teamNumber(forPlayer num: Int) -> Int {
        return num
    }
    
    func teamColor(forPlayer num: Int) -> UIColor {
        let team = teamNumber(forPlayer: num)
        return team == 0 ? color0 : color1
    }
    
    func teamColor(forTeam team: Int) -> UIColor {
        return team == 0 ? color0 : color1
    }
    
    func startingPosition(forPlayer num: Int) -> CGPoint {
        return startingPositions[num]!
    }
}
