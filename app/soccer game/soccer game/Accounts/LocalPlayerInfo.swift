//
//  LocalPlayerInfo.swift
//  soccer game
//
//  Created by rtoepfer on 4/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

struct LocalPlayerInfo {
    static var username : String?
    static var gamesPlayed : Int?
    static var goalsScored : Int?
    
    static var applicationToken : String?
}


struct CrudServiceLoginUser : Codable{
    
    let nickname: String?
    let gamesPlayed: String?
    let gamesWon: String?
    let goalsScored: String?
    
    private enum CodingKeys: String, CodingKey {
        case nickname
        case gamesPlayed = "gamesplayed"
        case gamesWon = "gameswon"
        case goalsScored = "goalsscored"
    }
}

struct CrudServiceLoginResponse : Codable {
    
    let profile : CrudServiceLoginUser?
    let token : String?
    
    private enum CodingKeys: String, CodingKey {
        case profile = "Profile"
        case token = "ApplicationToken"
    }
}

