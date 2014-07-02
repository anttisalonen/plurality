#include <cstring>
#include <fstream>
#include <iostream>
#include <vector>

#include <jsoncpp/json/json.h>

#include "GameObject.hpp"
#include "components/HelloComponent.hpp"

void runGame(const Json::Value& root)
{
	std::vector<GameObject> Objects;
	for(auto& jo : root["objects"]) {
		GameObject obj(jo["name"].asString());
		for(auto& jcomp : jo["components"]) {
			const std::string& type = jcomp["type"].asString();
			ComponentPtr comp;
			if(type == "HelloComponent") 
				comp = ComponentPtr(new HelloComponent());
			else
				std::cerr << "Invalid component type " << type << "!\n";

			if(comp) {
				obj.addComponent(comp);
				auto jvalnames = jcomp["values"].getMemberNames();
				for(auto& jvalname : jvalnames) {
					const std::string& valname = jvalname;
					const Json::Value& value = jcomp["values"][valname];
					if(value.isString())
						comp->addValue(valname, value.asString());
					else if(value.isIntegral())
						comp->addValue(valname, value.asInt());
					else
						std::cerr << "Invalid value type!\n";
				}
			}
		}
		Objects.push_back(obj);
	}

	for(auto& obj : Objects) {
		for(auto& comp : obj.getComponents())
			comp->Start();
	}

	Objects.clear();
}

int outputInterface(const char* filename)
{
	Json::Value root;
	root["components"] = Json::Value();

	ComponentPtr comp = ComponentPtr(new HelloComponent());
	Json::Value jcomp;
	jcomp["name"] = comp->getName();
	Json::Value jvalues;

	auto m = comp->getPossibleValues();
	for(const auto& mp : m) {
		jvalues[mp.first] = mp.second;
	}
	jcomp["values"] = jvalues;
	root["components"].append(jcomp);

	Json::StyledWriter writer;
	std::ofstream out(filename);
	out << writer.write(root);
	return 0;
}

int main(int argc, char** argv)
{
	if(argc == 3) {
		if(!strcmp(argv[1], "-o")) {
			outputInterface(argv[2]);
			exit(0);
		}
	}
	if(argc != 2) {
		std::cerr << "Usage: " << argv[0] << " <game JSON file>\n";
		exit(1);
	}
	std::string jsonFilename = argv[1];

	Json::Reader reader;
	Json::Value root;

	std::ifstream input(jsonFilename, std::ifstream::binary);
	bool parsingSuccessful = reader.parse(input, root, false);
	if (!parsingSuccessful) {
		throw std::runtime_error(reader.getFormatedErrorMessages());
	}

	runGame(root);

	return 0;
}
