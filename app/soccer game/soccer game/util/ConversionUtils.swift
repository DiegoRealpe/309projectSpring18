//
//  ByteConversions.swift
//  soccer game
//
//  Created by rtoepfer on 2/16/18.
//  Copyright © 2018 MG 6. All rights reserved.
//

import Foundation

func convertToFloat(_ bytes: [UInt8]) -> Float32{
    var f:Float32 = 0.0
    memcpy(&f, bytes, 4)
    return f
}

func convertToUInt8(_ float : Float32) -> [UInt8]{
    var float2 = float
    var bytes: [UInt8] = [0,0,0,0]
    memcpy(&bytes[0], &float2, 4)
    return bytes
}
