/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package hessian

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/alipay/sofa-mosn/pkg/log"
)

func Test_HessianCodecHeader(t *testing.T) {
	log.InitDefaultLogger("", log.INFO)

	// 这是应用层对象.
	strEchoBytes := []byte{0x4f, 0xbc, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x6c, 0x69, 0x70, 0x61, 0x79, 0x2e, 0x73, 0x6f, 0x66, 0x61, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x6f, 0x66, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x95, 0x0d, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x41, 0x70, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x17, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x73, 0x0d, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x41, 0x72, 0x67, 0x53, 0x69, 0x67, 0x73, 0x6f, 0x90, 0x04, 0x74, 0x65, 0x73, 0x74, 0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x49, 0x6e, 0x74, 0x53, 0x00, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x6c, 0x69, 0x70, 0x61, 0x79, 0x2e, 0x62, 0x65, 0x61, 0x6e, 0x2e, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x32, 0x3a, 0x31, 0x2e, 0x30, 0x4d, 0x11, 0x72, 0x70, 0x63, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x4d, 0x0d, 0x73, 0x6f, 0x66, 0x61, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x41, 0x70, 0x70, 0x04, 0x74, 0x65, 0x73, 0x74, 0x0c, 0x73, 0x6f, 0x66, 0x61, 0x50, 0x65, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x73, 0x00, 0x0b, 0x73, 0x79, 0x73, 0x50, 0x65, 0x6e, 0x41, 0x74, 0x74, 0x72, 0x73, 0x00, 0x0b, 0x73, 0x6f, 0x66, 0x61, 0x54, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x1e, 0x30, 0x61, 0x30, 0x66, 0x65, 0x38, 0x66, 0x38, 0x31, 0x35, 0x32, 0x32, 0x32, 0x30, 0x37, 0x35, 0x31, 0x30, 0x39, 0x32, 0x35, 0x31, 0x30, 0x30, 0x31, 0x31, 0x30, 0x35, 0x38, 0x32, 0x09, 0x73, 0x6f, 0x66, 0x61, 0x52, 0x70, 0x63, 0x49, 0x64, 0x01, 0x30, 0x0c, 0x73, 0x6f, 0x66, 0x61, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x70, 0x0d, 0x31, 0x30, 0x2e, 0x31, 0x35, 0x2e, 0x32, 0x33, 0x32, 0x2e, 0x32, 0x34, 0x38, 0x7a, 0x7a, 0x56, 0x74, 0x00, 0x07, 0x5b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x6e, 0x00, 0x7a}

	buffer2 := bytes.NewReader(strEchoBytes)
	decoder := NewDecoder(*buffer2, nil)

	decoder.RegisterType("com.alipay.sofa.rpc.core.request.SofaRequest", reflect.TypeOf(SofaRequest{}))
	obj, err := decoder.ReadObject()

	if err != nil {
		log.DefaultLogger.Errorf("read object err: %s", err)
	} else {
		if so, ok := obj.(reflect.Value); ok {
			u1 := so.Interface().(*SofaRequest) //

			log.DefaultLogger.Infof("obj,%+v", u1)
		}
	}
}

func Test_RequestIDCodec(t *testing.T) {

	log.InitDefaultLogger("", log.INFO)

	// 通信层对象的二进制
	strEchoBytes := []byte{0x4f, 0xba, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6f, 0x62, 0x61, 0x6f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x69, 0x6d, 0x70, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x91, 0x03, 0x63, 0x74, 0x78, 0x6f, 0x90, 0x4f, 0xc8, 0x39, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6f, 0x62, 0x61, 0x6f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x69, 0x6d, 0x70, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x24, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x92, 0x02, 0x69, 0x64, 0x06, 0x74, 0x68, 0x69, 0x73, 0x24, 0x30, 0x6f, 0x91, 0xe2, 0x4a, 0x00}

	buffer2 := bytes.NewReader(strEchoBytes)
	decoder := NewDecoder(*buffer2, nil)

	decoder.RegisterType("com.taobao.remoting.impl.ConnectionRequest", reflect.TypeOf(ConnectionRequest{}))
	decoder.RegisterType("com.taobao.remoting.impl.ConnectionRequest$RequestContext", reflect.TypeOf(RequestContext{}))
	obj, err := decoder.ReadObject()

	if err != nil {
		log.DefaultLogger.Errorf("read object err: %s", err)
	} else {
		if so, ok := obj.(reflect.Value); ok {
			u1 := so.Interface().(*ConnectionRequest)
			log.DefaultLogger.Infof("obj,%+v", u1)

			log.DefaultLogger.Infof(string(u1.Ctx.Id)) //GET REQUEST ID
		}
	}
}

func Test_ResponseIDCodec(t *testing.T) {

	log.InitDefaultLogger("", log.INFO)

	// 通信层响应对象的二进制
	strEchoBytes := []byte{0x4f, 0xbb, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6f, 0x62, 0x61, 0x6f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x69, 0x6d, 0x70, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x95, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x08, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x03, 0x63, 0x74, 0x78, 0x6f, 0x90, 0x13, 0x31, 0x30, 0x2e, 0x31, 0x35, 0x2e, 0x32, 0x33, 0x32, 0x2e, 0x32, 0x34, 0x38, 0x3a, 0x31, 0x32, 0x32, 0x39, 0x39, 0x90, 0x4e, 0x4e, 0x4f, 0xc8, 0x3b, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x6f, 0x62, 0x61, 0x6f, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x69, 0x6d, 0x70, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x24, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x92, 0x02, 0x69, 0x64, 0x06, 0x74, 0x68, 0x69, 0x73, 0x24, 0x30, 0x6f, 0x91, 0xe1, 0x4a, 0x00}

	buffer2 := bytes.NewReader(strEchoBytes)
	decoder := NewDecoder(*buffer2, nil)

	decoder.RegisterType("com.taobao.remoting.impl.ConnectionResponse", reflect.TypeOf(ConnectionResponse{}))
	decoder.RegisterType("com.taobao.remoting.impl.ConnectionResponse$ResponseContext", reflect.TypeOf(ResponseContext{}))
	obj, err := decoder.ReadObject()

	if err != nil {
		log.DefaultLogger.Errorf("read object err: %s", err)
	} else {
		if so, ok := obj.(reflect.Value); ok {
			u1 := so.Interface().(*ConnectionResponse)
			log.DefaultLogger.Infof("obj,%+v", u1)

			log.DefaultLogger.Infof(string(u1.Ctx.Id)) //GET RESPONSE ID
		}
	}
}
